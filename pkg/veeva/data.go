package veeva

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
)

func GetActivityEvents() ([]ActivityEvent, error) {
	db, err := sql.Open(settings.Luckie.DBDriver, conn)
	if err != nil {
		return nil, errors.Wrap(err, "Error from sql.Open() in GetActivityEvents() | driver: " + settings.Luckie.DBDriver + ", conn: " + conn)
	}
	defer db.Close()
	stmt, err := db.Prepare("select H_ACT_EVENT_ID, CID, SBL_CONTACT_ID, INTG_CONTACT_ID, START_TIME, END_TIME, AE_TYPE, AE_DESC, AE_SUBTYPE, INTG_AE_ID, TRIGGER_SRC_TYPE from H_ACT_EVENT")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	activityEvents := make([]ActivityEvent, 0)
	for rows.Next() {
		ae := ActivityEvent{}
		var H_ACT_EVENT_ID sql.NullInt64
		var CID, SBL_CONTACT_ID, INTG_CONTACT_ID, START_TIME, END_TIME, AE_TYPE, AE_DESC, AE_SUBTYPE, INTG_AE_ID, TRIGGER_SRC_TYPE sql.NullString
		var v interface{}
		var err error
		err = rows.Scan(&H_ACT_EVENT_ID, &CID, &SBL_CONTACT_ID, &INTG_CONTACT_ID, &START_TIME, &END_TIME, &AE_TYPE, &AE_DESC, &AE_SUBTYPE, &INTG_AE_ID, &TRIGGER_SRC_TYPE)
		if err != nil {
			return nil, err
		}
		if H_ACT_EVENT_ID.Valid {
			v, err = H_ACT_EVENT_ID.Value()
			if err != nil {
				return nil, err
			}
			ae.H_ACT_EVENT_ID = strings.TrimSpace(strconv.FormatInt(v.(int64),10))
		}
		if CID.Valid {
			v, err = CID.Value()
			if err != nil {
				return nil, err
			}
			ae.CID = strings.TrimSpace(v.(string))
		}
		if SBL_CONTACT_ID.Valid {
			v, err = SBL_CONTACT_ID.Value()
			if err != nil {
				return nil, err
			}
			ae.SBL_CONTACT_ID = strings.TrimSpace(v.(string))
		}
		if INTG_CONTACT_ID.Valid {
			v, err = INTG_CONTACT_ID.Value()
			if err != nil {
				return nil, err
			}
			ae.INTG_CONTACT_ID = strings.TrimSpace(v.(string))
		}
		if START_TIME.Valid {
			v, err = START_TIME.Value()
			if err != nil {
				return nil, err
			}
			ae.START_TIME = strings.TrimSpace(v.(string))
		}
		if END_TIME.Valid {
			v, err = END_TIME.Value()
			if err != nil {
				return nil, err
			}
			ae.END_TIME = strings.TrimSpace(v.(string))
		}
		if AE_TYPE.Valid {
			v, err = AE_TYPE.Value()
			if err != nil {
				return nil, err
			}
			ae.AE_TYPE = strings.TrimSpace(v.(string))
		}
		if AE_DESC.Valid {
			v, err = AE_DESC.Value()
			if err != nil {
				return nil, err
			}
			ae.AE_DESC = strings.TrimSpace(v.(string))
		}
		if AE_SUBTYPE.Valid {
			v, err = AE_SUBTYPE.Value()
			if err != nil {
				return nil, err
			}
			ae.AE_SUBTYPE = strings.TrimSpace(v.(string))
		}
		if INTG_AE_ID.Valid {
			v, err = INTG_AE_ID.Value()
			if err != nil {
				return nil, err
			}
			ae.INTG_AE_ID = strings.TrimSpace(v.(string))
		}
		if TRIGGER_SRC_TYPE.Valid {
			v, err = TRIGGER_SRC_TYPE.Value()
			if err != nil {
				return nil, err
			}
			ae.TRIGGER_SRC_TYPE = strings.TrimSpace(v.(string))
		}
		activityEvents = append(activityEvents, ae)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return activityEvents, nil
}

func CreateRecord(ae ActivityEvent) error {
	db, err := sql.Open(settings.Luckie.DBDriver, conn)
	if err != nil {
		return errors.Wrap(err, "Error from sql.Open() in CreateRecord")
	}
	defer db.Close()
	stmt, err := db.Prepare("insert into H_ACT_EVENT (CID, SBL_CONTACT_ID, INTG_CONTACT_ID, START_TIME, END_TIME, AE_TYPE, AE_DESC, AE_SUBTYPE, INTG_AE_ID) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ae.CID, ae.SBL_CONTACT_ID, ae.INTG_CONTACT_ID, ae.START_TIME, ae.END_TIME, ae.AE_TYPE, ae.AE_DESC, ae.AE_SUBTYPE, ae.INTG_AE_ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRecord(id string) error {
	db, err := sql.Open(settings.Luckie.DBDriver, conn)
	if err != nil {
		return errors.Wrap(err, "Error from sql.Open() in DeleteRecord()")
	}
	defer db.Close()
	stmt, err := db.Prepare("delete from H_ACT_EVENT where H_ACT_EVENT_ID = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func logHCPIDs(ids []id) error {
	db, err := sql.Open(settings.Luckie.DBDriver, conn)
	if err != nil {
		return errors.Wrap(err, "Error from sql.Open() in logHCPIDs()")
	}
	defer db.Close()
	for _, id := range ids {
		stmt, err := db.Prepare("insert into VeeVa_HCP_Info (CID, NewCID, SBL_CONTACT_ID, CreateDate) values ($1, $2, $3, $4)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(id.InitialCID, id.FinalCID, id.VeevaID, time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}

func logError(log Log) error {
	db, err := sql.Open(settings.Luckie.DBDriver, conn)
	if err != nil {
		return errors.Wrap(err, "Error from sql.Open() in logError()")
	}
	defer db.Close()
	stmt, err := db.Prepare("insert into VeeVa_SOAP_Log (CID, RequestType, Request, Response, Error, DateCreated, FinalCID, SBL_CONTACT_ID, DateSent, DateReceived, Status) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(log.CID, log.RequestType, log.Request, log.Response, log.Error, log.DateCreated, log.FinalCID, log.SBL_CONTACT_ID, log.DateSent, log.DateReceived, log.Status)
	if err != nil {
		return err
	}
	return nil
}