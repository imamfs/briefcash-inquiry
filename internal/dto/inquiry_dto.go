package dto

type InquiryResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Source  string      `json:"source"`
	Data    InquiryData `json:"data"`
}

type InquiryRequest struct {
	CompanyId          string `json:"company_id"`
	BeneficiaryAccount string `json:"beneficary_account"`
	PartnerReferenceNo string `json:"partner_reference_no"`
	BankCode           string `json:"bank_code"`
	Type               string `json:"type"` // bifast, online
}

type InquiryData struct {
	BeneficiaryAccount string `json:"beneficary_account"`
	BankCode           string `json:"bank_code"`
	BeneficiaryName    string `json:"beneficiary_name"`
}

type BCAInternalInquiryRequest struct {
	PartnerReferenceNo   string `json:"partnerReferenceNo"`
	BeneficiaryAccountNo string `json:"beneficiaryAccountNo"`
}

type BCAExternalInquiryRequest struct {
	BeneficiaryBankCode  string            `json:"beneficiaryBankCode"`
	BeneficiaryAccountNo string            `json:"beneficiaryAccountNo"`
	PartnerReferenceNo   string            `json:"partnerReferenceNo"`
	AdditionalInfo       BCAAdditionalInfo `json:"additionalInfo"`
}

type BCAAdditionalInfo struct {
	InquiryService string `json:"inquiryService"`
}

type BRIInternalInquiryRequest struct {
	BeneficiaryAccountNo string            `json:"beneficiaryAccountNo"`
	AdditionalInfo       map[string]string `json:"additionalInfo"`
}

type BRIExternalInquiryRequest struct {
	BeneficiaryBankCode  string            `json:"beneficiaryBankCode"`
	BeneficiaryAccountNo string            `json:"beneficiaryAccountNo"`
	AdditionalInfo       map[string]string `json:"additionalInfo"`
}

type CimbInternalInquiryRequest struct {
	PartnerReferenceNo   string            `json:"partnerReferenceNo"`
	BeneficiaryAccountNo string            `json:"beneficiaryAccountNo"`
	AdditionalInfo       map[string]string `json:"additionalInfo"`
}

type CimbExternalInquiryRequest struct {
	PartnerReferenceNo   string            `json:"partnerReferenceNo"`
	BeneficiaryBankCode  string            `json:"beneficiaryBankCode"`
	BeneficiaryAccountNo string            `json:"beneficiaryAccountNo"`
	AdditionalInfo       map[string]string `json:"additionalInfo"`
}

type PermataInquiryHeaderRequest struct {
	RequestTimeStamp string `json:"RequestTimeStamp"`
	CustReffID       string `json:"CustReffID"`
}

type PermataInternalInquiryBodyRequest struct {
	AccountNumber string `json:"AccountNumber"`
}

type PermataInternalInquiryRequest struct {
	MessageHeader PermataInquiryHeaderRequest       `json:"MsgRqHdr"`
	MessageBody   PermataInternalInquiryBodyRequest `json:"InqInfo"`
}

type PermataExternalInquiryBodyRequest struct {
	ToAccount string `json:"ToAccount"`
	BankId    string `json:"BankId"`
	BankName  string `json:"BankName"`
}

type PermataExternalInquiryRequest struct {
	MessageHeader PermataInquiryHeaderRequest       `json:"MsgRqHdr"`
	MessageBody   PermataExternalInquiryBodyRequest `json:"XferInfo"`
}

type BCAInquiryResponse struct {
	ResponseCode           string `json:"responseCode"`
	ResponseMessage        string `json:"responseMessage"`
	ReferenceNo            string `json:"referenceNo"`
	PartnerReferenceNo     string `json:"partnerReferenceNo"`
	BeneficiaryAccountName string `json:"beneficiaryAccountName"`
	BeneficiaryAccountNo   string `json:"beneficiaryAccountNo"`
	BeneficiaryBankCode    string `json:"beneficiaryBankCode"`
}

type BRIInternalInquiryResponse struct {
	ResponseCode             string            `json:"responseCode"`
	ResponseMessage          string            `json:"responseMessage"`
	ReferenceNo              string            `json:"referenceNo"`
	BeneficiaryAccountNo     string            `json:"beneficiaryAccountNo"`
	BeneficiaryAccountName   string            `json:"beneficiaryAccountName"`
	BeneficiaryAccountStatus string            `json:"beneficiaryAccountStatus"`
	BeneficiaryAccountType   string            `json:"beneficiaryAccountType"`
	Currency                 string            `json:"currency"`
	AdditionalInfo           map[string]string `json:"additionalInfo"`
}

type BRIExternalInquiryResponse struct {
	ResponseCode           string            `json:"responseCode"`
	ResponseMessage        string            `json:"responseMessage"`
	ReferenceNo            string            `json:"referenceNo"`
	BeneficiaryAccountName string            `json:"beneficiaryAccountName"`
	BeneficiaryAccountNo   string            `json:"beneficiaryAccountNo"`
	BeneficiaryBankCode    string            `json:"beneficiaryBankCode"`
	BeneficiaryBankName    string            `json:"beneficiaryBankName"`
	Currency               string            `json:"currency"`
	AdditionalInfo         map[string]string `json:"additionalInfo"`
}

type BRIErrorResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}

type CimbInternalInquiryResponse struct {
	ResponseCode             string            `json:"responseCode"`
	ResponseMessage          string            `json:"responseMessage"`
	PartnerReferenceNo       string            `json:"partnerReferenceNo"`
	BeneficiaryAccountName   string            `json:"beneficiaryAccountName"`
	BeneficiaryAccountNo     string            `json:"beneficiaryAccountNo"`
	BeneficiaryAccountStatus string            `json:"beneficiaryAccountStatus"`
	BeneficiaryAccountType   string            `json:"beneficiaryAccountType"`
	Currency                 string            `json:"currency"`
	AdditionalInfo           map[string]string `json:"additionalInfo"`
}

type CimbExternalInquiryResponse struct {
	ResponseCode           string            `json:"responseCode"`
	ResponseMessage        string            `json:"responseMessage"`
	PartnerReferenceNo     string            `json:"partnerReferenceNo"`
	BeneficiaryAccountName string            `json:"beneficiaryAccountName"`
	BeneficiaryAccountNo   string            `json:"beneficiaryAccountNo"`
	BeneficiaryBankCode    string            `json:"beneficiaryBankCode"`
	Currency               string            `json:"currency"`
	AdditionalInfo         map[string]string `json:"additionalInfo"`
}

type PermataInquiryHeaderResponse struct {
	ResponseTimestamp string `json:"ResponseTimestamp"`
	CustReffID        string `json:"CustReffID"`
	StatusCode        string `json:"StatusCode"`
	StatusDesc        string `json:"StatusDesc"`
}

type PermataInternalInquiryBodyResponse struct {
	AccountNumber string `json:"AccountNumber"`
	AccountName   string `json:"AccountName"`
}

type PermataInternalInquiryResponse struct {
	MessageHeader PermataInquiryHeaderResponse       `json:"MsgRsHdr"`
	MessageBody   PermataInternalInquiryBodyResponse `json:"InqInfo"`
}

type PermataExternalInquiryBodyResponse struct {
	ToAccount         string `json:"ToAccount"`
	ToAccountFullName string `json:"ToAccountFullName"`
	BankId            string `json:"BankId"`
	BankName          string `json:"BankName"`
}

type PermataExternalInquiryResponse struct {
	MessageHeader PermataInquiryHeaderResponse       `json:"MsgRsHdr"`
	MessageBody   PermataExternalInquiryBodyResponse `json:"InqInfo"`
}

type SNAPAccessToken struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	AccessToken     string `json:"accessToken"`
	TokenType       string `json:"tokenType"`
	ExpiresIn       int16  `json:"expiresIn"`
}
