package veeva

var headerSOAP string
var cdiSOAP string
var conVerifySOAP string
var putActEventsSOAP string

func initTemplates() {

cdiSOAP = "{{define \"cdi\"}}" +
"<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
"<soapenv:Envelope xmlns:cus=\"http://ciis.gsk.com/CustomerService\" xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
  "<soapenv:Header>" +
     "<wsse:Security soapenv:mustUnderstand=\"0\" xmlns:wsse=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd\">" +
        "<wsse:UsernameToken wsu:Id=\"UsernameToken-1\" xmlns:wsu=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd\">" +
           "<wsse:Username>" + settings.GSK.Username + "</wsse:Username>" +
           "<wsse:Password Type=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText\">" + settings.GSK.Password + "</wsse:Password>" +
           "<wsse:Nonce EncodingType=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary\">" + settings.GSK.CdiNonce + "</wsse:Nonce>" +
           "<wsu:Created>2016-09-27T02:46:25.831Z</wsu:Created>" +
        "</wsse:UsernameToken>" +
     "</wsse:Security>" +
     "<cus:TransactionHeader>" +
       "<cus:BusinessTransactionID>0</cus:BusinessTransactionID>" +
       "<cus:SourceEAN>0</cus:SourceEAN>" +
     "</cus:TransactionHeader>" +
  "</soapenv:Header>" +
  "<soapenv:Body>" +
    "<cus:SearchAltIdsRequest>" +
      "<cus:request>" +
        "<cus:CID>{{ .InitialCID }}</cus:CID>" +
      "</cus:request>" +
    "</cus:SearchAltIdsRequest>" +
  "</soapenv:Body>" +
"</soapenv:Envelope>" +
"{{end}}"

headerSOAP = "{{define \"header\"}}" +
"<soapenv:Header>" +
  "<wsse:Security soapenv:mustUnderstand=\"1\" xmlns:wsse=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd\" xmlns:wsu=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd\">" +
    "<wsse:UsernameToken wsu:Id=\"UsernameToken-BF262DFEA6FD34A40E14279111853981\">" +
      "<wsse:Username>" + settings.GSK.Username + "</wsse:Username>" +
      "<wsse:Password Type=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText\">" + settings.GSK.Password + "</wsse:Password>" +
      "<wsse:Nonce EncodingType=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary\">" + settings.GSK.IdsNonce + "</wsse:Nonce>" +
      "<wsu:Created>2016-10-17T16:00:00Z</wsu:Created>" +
    "</wsse:UsernameToken>" +
  "</wsse:Security>" +
"</soapenv:Header>" +
"{{end}}"

conVerifySOAP = "{{define \"conVerify\"}}" +
"<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
"<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:ccit=\"http://ccit.gsk.org/\">" +
  "{{template \"header\"}}" +
  "<soapenv:Body>" +
    "<ccit:con_verify>" +
      "<ccit:buId>US</ccit:buId>" +
      "<ccit:maxRows>1</ccit:maxRows>" +
      "<ccit:legalId>{{ .FinalCID }}</ccit:legalId>" +
    "</ccit:con_verify>" +
  "</soapenv:Body>" +
"</soapenv:Envelope>" +
"{{end}}"

putActEventsSOAP = "{{define \"putActEvents\"}}" +
"<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
"<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:ccit=\"http://ccit.gsk.org/\">" +
  "{{template \"header\"}}" +
  "<soapenv:Body>" +
    "<ccit:PutActEvents>" +
      "<ccit:bu_id>US</ccit:bu_id>" +
      "<ccit:actEvents>" +
        "<ccit:H_ACT_EVENT>" +
          "<ccit:BU_ID>US</ccit:BU_ID>" +
          "<ccit:SBL_CONTACT_ID>{{ .SBL_CONTACT_ID }}</ccit:SBL_CONTACT_ID>" +
          "<ccit:INTG_CONTACT_ID>{{ .INTG_CONTACT_ID }}</ccit:INTG_CONTACT_ID>" +
          "<ccit:START_TIME>{{ .START_TIME }}</ccit:START_TIME>" +
          "<ccit:END_TIME>{{ .END_TIME }}</ccit:END_TIME>" +
          "<ccit:AE_TYPE>{{ .AE_TYPE }}</ccit:AE_TYPE>" +
          "<ccit:DATASOURCE_CD>LUCKIE</ccit:DATASOURCE_CD>" +
          "<ccit:AE_STATUS>SUBMITTED</ccit:AE_STATUS>" +
          "<ccit:AE_DESC>{{ .AE_DESC }}</ccit:AE_DESC>" +
          "<ccit:AE_SUBTYPE>{{ .AE_SUBTYPE }}</ccit:AE_SUBTYPE>" +
          "<ccit:INTG_AE_ID>{{ .INTG_AE_ID }}</ccit:INTG_AE_ID>" +
          "<ccit:OWNER>Luckie</ccit:OWNER>" +
          "<ccit:TRIGGER_SRC_TYPE>{{ .TRIGGER_SRC_TYPE }}</ccit:TRIGGER_SRC_TYPE>" +
        "</ccit:H_ACT_EVENT>" +
      "</ccit:actEvents>" +
    "</ccit:PutActEvents>" +
  "</soapenv:Body>" +
"</soapenv:Envelope>" +
"{{end}}"

}
