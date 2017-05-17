package veeva

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

type id struct {
	InitialCID string
	FinalCID   string
	VeevaID    string
}

type ActivityEvent struct {
	H_ACT_EVENT_ID   string
	CID              string
	SBL_CONTACT_ID   string
	INTG_CONTACT_ID  string
	START_TIME       string
	END_TIME         string
	AE_TYPE          string
	AE_DESC          string
	AE_SUBTYPE       string
	INTG_AE_ID       string
	TRIGGER_SRC_TYPE string
}

type Log struct {
	CID		string
	RequestType 	string
	Request		string
	Response 	string
	Error 		string
	DateCreated 	time.Time
	FinalCID	string
	SBL_CONTACT_ID	string
	DateSent	time.Time
	DateReceived	time.Time
	Status		string
}

type uniqueCIDs map[string]bool

func (cids uniqueCIDs) add(cid string) bool {
	cid = strings.TrimSpace(cid)
	if _, ok := cids[cid]; ok {
		return false
	}
	cids[cid] = true
	return true
}

func StartClock() error {
	hour := os.Getenv("HOUR")
	minute := os.Getenv("MINUTE")
	if hour == "" || minute == "" {
		return errors.New("HOUR and/or MINUTE environmental variables not set.")
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	c := cron.New()
	c.AddFunc("0 0 " + hour + " " + minute + " * *", func() { Start("") })
	c.Start()
	defer c.Stop()
	wg.Wait()
	return nil
}

func Start(awsProfile string) error {
	LoadConfig(awsProfile)
	initTemplates()
	activityEvents, err := GetActivityEvents()
	if err != nil {
		return errors.Wrap(err, "Error in GetActivityEvents()")
	}
	ids := getActivityEventCIDs(activityEvents)
	ids, err = verifyCIDs(ids)
	if err != nil {
		return errors.Wrap(err, "Error in verifyCIDs()")
	}
	ids, err = getVeevaIDs(ids)
	if err != nil {
		return errors.Wrap(err, "Error in getVeevaIDs()")
	}
	err = logHCPIDs(ids)
	if err != nil {
		return errors.Wrap(err, "Error in logHCPIDs()")
	}
	accepted, _ := updateActivityEvents(activityEvents, ids)  // _ is rejected
	err = uploadActivityEvents(accepted)
	if err != nil {
		return errors.Wrap(err, "Error in uploadActivityEvents()")
	}
	return nil
}

/*

COMMENTED OUT FUNCTIONS ARE BEGINING OF WORK-IN-PROGRESS
TO IMPLEMENT RATE-LIMITED CONCURRENCY TO ENHANCE PERFORMANCE.

func getCIDs(inputIDs []id) ([]id, error) {
	var wg sync.WaitGroup
	wg.Add(len(inputIDs))

	outputIDs := []id{}
	for _, id := range inputIDs {
		go getCID(id, &outputIDs, &wg)
	}
	wg.Wait()
	return outputIDs, nil
}

func getCID(thisID id, outputIDs *[]id, wg *sync.WaitGroup) {

	defer wg.Done()
	t := template.Must(template.New("cdi").Parse(cdiSOAP))
	var b bytes.Buffer
	err := t.ExecuteTemplate(&b, "cdi", thisID)
	if err != nil {
		//errors.New("getCID() for initial CID: " + id + " | ExecuteTemplate() | " + err.Error())
	}
	if len(b.Bytes()) == 0 {
		//errors.New("getCID() for initial CID: " + id + " | bytes.Buffer b received zero bytes from ExecuteTemplate()")
	}
	soap, err := callWebService(cdiURL, b.Bytes(), cdiAuthToken)
	finalCID, err := ValueFromXML(soap, "CID")
	thisID.FinalCID = finalCID
	*outputIDs = append(*outputIDs, thisID)
}
*/

func getActivityEventCIDs(activityEvents []ActivityEvent) []id {
	set := make(uniqueCIDs)
	for _, ae := range activityEvents {
		_ = set.add(ae.CID)
	}
	ids := []id{}
	for k := range set {
		i := id{}
		i.InitialCID = k
		ids = append(ids, i)
	}
	return ids
}

func verifyCIDs(inputIDs []id) ([]id, error) {
	outputIDs := []id{}
	for _, id := range inputIDs {
		t := template.Must(template.New("cdi").Parse(cdiSOAP))
		var b bytes.Buffer
		err := t.ExecuteTemplate(&b, "cdi", id)
		if err != nil {
			return nil, errors.Wrap(err, "verifyCIDs() | ExecuteTemplate()")
		}
		if len(b.Bytes()) == 0 {
			return nil, errors.New("verifyCIDs() | bytes.Buffer b received zero bytes from ExecuteTemplate()")
		}
		soap, err := callWebService(settings.GSK.CdiURL, b.Bytes(), settings.GSK.CdiAuthToken)
		finalCID, err := ValueFromXML(soap, "CID")
		if err != nil || finalCID == "" {
			log := Log{
				CID: id.InitialCID,
				RequestType: "cdi",
				Request: b.String(),
				Response: string(soap),
				DateCreated: time.Now(),
				DateSent: time.Now(),
				DateReceived: time.Now(),
			}
			_ = logError(log)
		}
		id.FinalCID = finalCID
		outputIDs = append(outputIDs, id)
		// Needs to be removed after everything is working.
		fmt.Println("================================================================================")
		fmt.Println("====== SOAP sent to CDI for CID: " + id.InitialCID)
		fmt.Println("================================================================================")
		fmt.Println(b.String())
		fmt.Println("================================================================================")
		fmt.Println("====== SOAP received from CDI for CID: " + id.InitialCID)
		fmt.Println("================================================================================")
		fmt.Println(string(soap))
		fmt.Println("================================================================================")
		fmt.Println("Initial CID: " + id.InitialCID + " => Final CID: " + id.FinalCID)
		fmt.Println("================================================================================")
	}
	return outputIDs, nil
}

func getVeevaIDs(inputIDs []id) ([]id, error) {
	outputIDs := []id{}
	for _, id := range inputIDs {
		if id.FinalCID != "" {
			t := template.Must(template.New("conVerify").Parse(conVerifySOAP))
			t = template.Must(t.New("header").Parse(headerSOAP))
			var b bytes.Buffer
			err := t.ExecuteTemplate(&b, "conVerify", id)
			if err != nil {
				return nil, errors.Wrap(err, "getVeevaIDs() | ExecuteTemplate()")
			}
			if len(b.Bytes()) == 0 {
				return nil, errors.New("getVeevaIDs() | bytes.Buffer b received zero bytes from ExecuteTemplate()")
			}
			soap, err := callWebService(settings.GSK.IdsURL, b.Bytes(), "")
			veevaID, err := ValueFromXML(soap, "rowId")
			if err != nil || veevaID == "" {
				log := Log{
					CID: id.InitialCID,
					RequestType: "conVerify",
					Request: b.String(),
					Response: string(soap),
					DateCreated: time.Now(),
					FinalCID: id.FinalCID,
					DateSent: time.Now(),
					DateReceived: time.Now(),
				}
				_ = logError(log)
			}
			id.VeevaID = veevaID
			outputIDs = append(outputIDs, id)

			// Needs to be removed after everything is working.
			fmt.Println("================================================================================")
			fmt.Println("====== SOAP sent to IDS con_verify for CID: " + id.FinalCID)
			fmt.Println("================================================================================")
			fmt.Println(b.String())
			fmt.Println("================================================================================")
			fmt.Println("====== SOAP received from IDS con_verify for CID: " + id.FinalCID)
			fmt.Println("================================================================================")
			fmt.Println(string(soap))
			fmt.Println("================================================================================")
			fmt.Println("====== CID: " + id.FinalCID + " => Veeva ID: " + id.VeevaID)
			fmt.Println("================================================================================")
		}
	}
	return outputIDs, nil
}

func updateActivityEvents(activityEvents []ActivityEvent, ids []id) ([]ActivityEvent, []ActivityEvent) {

	accepted := []ActivityEvent{}
	rejected := []ActivityEvent{}
	for _, ae := range activityEvents {
		for _, i := range ids {
			if strings.TrimSpace(ae.CID) == strings.TrimSpace(i.InitialCID) {
				if i.VeevaID != "" {
					ae.SBL_CONTACT_ID = i.VeevaID
					accepted = append(accepted, ae)
				} else {
					rejected = append(rejected, ae)
				}
			}
		}
	}
	return accepted, rejected
}

func uploadActivityEvents(activityEvents []ActivityEvent) error {
	for _, ae := range activityEvents {
		if ae.SBL_CONTACT_ID == "" {
			fmt.Println("====== FAILURE PutActEvents for ID: " + ae.H_ACT_EVENT_ID + " has no Veeva ID.")
			return errors.New("uploadActivityEvents() | ActivityEvent ID: " + ae.H_ACT_EVENT_ID + " has no Veeva ID.")
		}
		t := template.Must(template.New("putActEvents").Parse(putActEventsSOAP))
		t = template.Must(t.New("header").Parse(headerSOAP))
		var b bytes.Buffer
		err := t.ExecuteTemplate(&b, "putActEvents", ae)
		if err != nil {
			// fmt.Println("====== FAILURE PutActEvents for ID: " + ae.H_ACT_EVENT_ID + ", CID: " + ae.CID + ", Veeva ID: " + ae.SBL_CONTACT_ID + ", INTG_AE_ID: " + ae.INTG_AE_ID + " | uploadActivityEvents() | ExecuteTemplate() | " + err.Error())
			return errors.Wrap(err, "uploadActivityEvents() | ExecuteTemplate()")
		}
		soap, err := callWebService(settings.GSK.IdsURL, b.Bytes(), "")
		if err != nil {
			log := Log{
				RequestType: "putActEvents",
				Request: b.String(),
				Response: string(soap),
				DateCreated: time.Now(),
				FinalCID: ae.CID,
				SBL_CONTACT_ID: ae.SBL_CONTACT_ID,
				DateSent: time.Now(),
				DateReceived: time.Now(),
			}
			_ = logError(log)
		}
		fmt.Println("====== SUCCESS PutActEvents for ID: " + ae.H_ACT_EVENT_ID + ", CID: " + ae.CID + ", Veeva ID: " + ae.SBL_CONTACT_ID + ", INTG_AE_ID: " + ae.INTG_AE_ID)

		err = DeleteRecord(ae.H_ACT_EVENT_ID)
		if err != nil {
			fmt.Println("====== FAILURE PutActEvents ID: " + ae.H_ACT_EVENT_ID + " uploaded successfully, but did not delete record from table.")
		}
		// Needs to be removed after everything is working.
		fmt.Println("================================================================================")
		fmt.Println("====== SOAP sent to IDS PutActEvents for CID: " + ae.CID + ", Veeva ID: " + ae.SBL_CONTACT_ID + ", INTG_AE_ID: " + ae.INTG_AE_ID)
		fmt.Println("================================================================================")
		fmt.Println(b.String())
		fmt.Println("================================================================================")
		fmt.Println("====== SOAP received from IDS PutActEvents for CID: " + ae.CID + ", Veeva ID: " + ae.SBL_CONTACT_ID + ", INTG_AE_ID: " + ae.INTG_AE_ID)
		fmt.Println("================================================================================")
		fmt.Println(string(soap))

	}
	return nil
}

func callWebService(url string, payload []byte, authToken string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.New("http.NewRequest() | " + err.Error())
	}
	req.Header.Set("Content-Type", "application/soap+xml")
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("client.Do() | " + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("ioutil.ReadAll() | " + err.Error())
	}
	return body, nil
}

func ValueFromXML(b []byte, targetElement string) (string, error) {
	buf := bytes.NewBuffer(b)
	decoder := xml.NewDecoder(buf)
	var value string
	var inElement string
	for {
		// Read tokens from the XML document in a stream.
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		if se, ok := token.(xml.StartElement); ok {
			inElement = se.Name.Local
			if inElement == targetElement {
				targetToken, _ := decoder.Token()
				if v, ok := targetToken.(xml.CharData); ok {
					value = string(v)
				}
			}
		}

	}
	return value, nil
}
