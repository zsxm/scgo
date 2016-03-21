package soap

//WSDL 根元素
type Definitions struct {
	Xmlns           xmlns
	Name            string //MessageServiceJaxbImplService
	TargetNamespace string //http://ws.service.creditease.com/
	Xmlns           string "http://schemas.xmlsoap.org/wsdl/"
}
type xmlns struct {
	Tns       string "" //http://ws.service.creditease.com/
	Wsu       string "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	Wsp       string "http://www.w3.org/ns/ws-policy"
	Wsp1_2    string "http://schemas.xmlsoap.org/ws/2004/09/policy"
	Wsdlsoap  string "http://schemas.xmlsoap.org/wsdl/soap/"
	Soap12    string "http://www.w3.org/2003/05/soap-envelope"
	Xsd       string "http://www.w3.org/2001/XMLSchema"
	Soapenc11 string "http://schemas.xmlsoap.org/soap/encoding/"
	Soapenc12 string "http://www.w3.org/2003/05/soap-encoding"
	Soap11    string "http://schemas.xmlsoap.org/soap/envelope/"
	Soap      string "http://schemas.xmlsoap.org/wsdl/soap/"
	Wsdl      string "http://schemas.xmlsoap.org/wsdl/"
	Wsam      string "http://www.w3.org/2007/05/addressing/metadata"
}
type types struct {
}
