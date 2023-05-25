package tukpdq

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukhttp"
	"github.com/ipthomas/tukutil"
)

type PDQQuery struct {
	Server_Mode      string           `json:",omitempty"`
	Server_URL       string           `json:",omitempty"`
	CGL_X_Api_Key    string           `json:",omitempty"`
	CGL_X_Api_Secret string           `json:",omitempty"`
	NHS_ID           string           `json:",omitempty"`
	NHS_OID          string           `json:",omitempty"`
	MRN_ID           string           `json:",omitempty"`
	MRN_OID          string           `json:",omitempty"`
	REG_ID           string           `json:",omitempty"`
	REG_OID          string           `json:",omitempty"`
	GivenName        string           `json:"givenname"`
	FamilyName       string           `json:"familyname"`
	BirthDate        string           `json:"birthdate"`
	Gender           string           `json:"gender"`
	Zip              string           `json:"zip"`
	Street           string           `json:"street"`
	Town             string           `json:"town"`
	City             string           `json:"city"`
	Country          string           `json:"country"`
	Timeout          int64            `json:",omitempty"`
	Cache            bool             `json:",omitempty"`
	Used_PID         string           `json:",omitempty"`
	Used_PID_OID     string           `json:",omitempty"`
	Request          []byte           `json:",omitempty"`
	Response         []byte           `json:",omitempty"`
	StatusCode       int              `json:",omitempty"`
	Count            int              `json:",omitempty"`
	PDQv3Response    *PDQv3Response   `json:",omitempty"`
	PIXv3Response    *PIXv3Response   `json:",omitempty"`
	PIXmResponse     *PIXmResponse    `json:",omitempty"`
	Patients         *[]TUKPatient    `json:",omitempty"`
	CGLUserResponse  *CGLUserResponse `json:",omitempty"`
}
type Delphi struct {
	Data struct {
		LocalIdentifier int    `json:"LocalIdentifier,omitempty"`
		Status          string `json:"Status,omitempty"`
		Title           string `json:"Title,omitempty"`
		Forename        string `json:"Forename,omitempty"`
		Surname         string `json:"Surname,omitempty"`
		GenderAtBirth   string `json:"GenderAtBirth,omitempty"`
		DateOfBirth     string `json:"DateOfBirth,omitempty"`
		Address         struct {
			LocalIdentifier int    `json:"LocalIdentifier,omitempty"`
			AddressLine1    string `json:"AddressLine1,omitempty"`
			AddressLine2    string `json:"AddressLine2,omitempty"`
			AddressLine3    string `json:"AddressLine3,omitempty"`
			AddressLine4    string `json:"AddressLine4,omitempty"`
			PostCode1       string `json:"PostCode1,omitempty"`
			PostCode2       string `json:"PostCode2,omitempty"`
		} `json:"Address,omitempty"`
		Keyworker               string `json:"Keyworker,omitempty"`
		LastAttendedAppointment string `json:"LastAttendedAppointment,omitempty"`
		DrugScreening           []any  `json:"DrugScreening,omitempty"`
		Prescriptions           []any  `json:"Prescriptions,omitempty"`
		Risks                   []any  `json:"Risks,omitempty"`
		Careplans               []struct {
			LocalIdentifier          int    `json:"LocalIdentifier,omitempty"`
			AlcohoUse                bool   `json:"AlcohoUse,omitempty"`
			DrugUse                  bool   `json:"DrugUse,omitempty"`
			EffectsOfAlcoholAndDrugs bool   `json:"EffectsOfAlcoholAndDrugs,omitempty"`
			PreventingRelapse        bool   `json:"PreventingRelapse,omitempty"`
			PreventingOverdose       bool   `json:"PreventingOverdose,omitempty"`
			PersonalCare             bool   `json:"PersonalCare,omitempty"`
			FindingThingsIEnjoy      bool   `json:"FindingThingsIEnjoy,omitempty"`
			ManagingMoney            bool   `json:"ManagingMoney,omitempty"`
			SupportForMyChildren     bool   `json:"SupportForMyChildren,omitempty"`
			EducationOrTraining      bool   `json:"EducationOrTraining,omitempty"`
			Other                    bool   `json:"Other,omitempty"`
			AlcoholDrugUse           bool   `json:"AlcoholDrugUse,omitempty"`
			ManagingCravings         bool   `json:"ManagingCravings,omitempty"`
			MentalEmotionalHealth    bool   `json:"MentalEmotionalHealth,omitempty"`
			AccommodationHousing     bool   `json:"AccommodationHousing,omitempty"`
			LegalProblems            bool   `json:"LegalProblems,omitempty"`
			ParentingHelpSupport     bool   `json:"ParentingHelpSupport,omitempty"`
			PhyscialHealth           bool   `json:"PhyscialHealth,omitempty"`
			ImmediateProblem         string `json:"ImmediateProblem,omitempty"`
			LongTermGoal             string `json:"LongTermGoal,omitempty"`
			StepsToAchievingGoal     string `json:"StepsToAchievingGoal,omitempty"`
			HowDidItGo               string `json:"HowDidItGo,omitempty"`
			NextStepForGoal          string `json:"NextStepForGoal,omitempty"`
			CommunityDetox           bool   `json:"CommunityDetox,omitempty"`
			InpatientDetox           bool   `json:"InpatientDetox,omitempty"`
			OverdoseInformation      bool   `json:"OverdoseInformation,omitempty"`
			NutritionalAdvice        bool   `json:"NutritionalAdvice,omitempty"`
			HepCScreening            bool   `json:"HepCScreening,omitempty"`
			HepAAndBVaccination      bool   `json:"HepAAndBVaccination,omitempty"`
			GroupWork                bool   `json:"GroupWork,omitempty"`
			OneToOneSupport          bool   `json:"OneToOneSupport,omitempty"`
			SupportWorker            bool   `json:"SupportWorker,omitempty"`
			PrescribedMedication     bool   `json:"PrescribedMedication,omitempty"`
			Stabilisation            bool   `json:"Stabilisation,omitempty"`
			MedicalReview            bool   `json:"MedicalReview,omitempty"`
			OtherClinical            bool   `json:"OtherClinical,omitempty"`
			CarePlanGivenToClient    string `json:"CarePlanGivenToClient,omitempty"`
			CareplanStartDate        string `json:"CareplanStartDate,omitempty"`
			CarePlanReviewDate       string `json:"CarePlanReviewDate,omitempty"`
		} `json:"Careplans,omitempty"`
		Discharge any `json:"Discharge,omitempty"`
	} `json:"Data,omitempty"`
}

type CGLUserResponse struct {
	Data struct {
		Client struct {
			BasicDetails struct {
				Address struct {
					AddressLine1 string `json:"addressLine1,omitempty"`
					AddressLine2 string `json:"addressLine2,omitempty"`
					AddressLine3 string `json:"addressLine3,omitempty"`
					AddressLine4 string `json:"addressLine4,omitempty"`
					AddressLine5 string `json:"addressLine5,omitempty"`
					PostCode     string `json:"postCode,omitempty"`
				} `json:"address,omitempty"`
				BirthDate                    string `json:"birthDate,omitempty"`
				Disability                   string `json:"disability,omitempty"`
				LastEngagementByCGLDate      string `json:"lastEngagementByCGLDate,omitempty"`
				LastFaceToFaceEngagementDate string `json:"lastFaceToFaceEngagementDate,omitempty"`
				LocalIdentifier              int    `json:"localIdentifier,omitempty"`
				Name                         struct {
					Family string `json:"family,omitempty"`
					Given  string `json:"given,omitempty"`
				} `json:"name,omitempty"`
				NextCGLAppointmentDate string `json:"nextCGLAppointmentDate,omitempty"`
				NhsNumber              string `json:"nhsNumber,omitempty"`
				SexAtBirth             string `json:"sexAtBirth,omitempty"`
			} `json:"basicDetails,omitempty"`
			BbvInformation struct {
				BbvTested        string `json:"bbvTested,omitempty"`
				HepCLastTestDate string `json:"hepCLastTestDate,omitempty"`
				HepCResult       string `json:"hepCResult,omitempty"`
				HivPositive      string `json:"hivPositive,omitempty"`
			} `json:"bbvInformation,omitempty"`
			DrugTestResults struct {
				DrugTestDate          string `json:"drugTestDate,omitempty"`
				DrugTestSample        string `json:"drugTestSample,omitempty"`
				DrugTestStatus        string `json:"drugTestStatus,omitempty"`
				InstantOrConfirmation string `json:"instantOrConfirmation,omitempty"`
				Results               struct {
					Amphetamine     string `json:"amphetamine,omitempty"`
					Benzodiazepine  string `json:"benzodiazepine,omitempty"`
					Buprenorphine   string `json:"buprenorphine,omitempty"`
					Cannabis        string `json:"cannabis,omitempty"`
					Cocaine         string `json:"cocaine,omitempty"`
					Eddp            string `json:"eddp,omitempty"`
					Fentanyl        string `json:"fentanyl,omitempty"`
					Ketamine        string `json:"ketamine,omitempty"`
					Methadone       string `json:"methadone,omitempty"`
					Methamphetamine string `json:"methamphetamine,omitempty"`
					Morphine        string `json:"morphine,omitempty"`
					Opiates         string `json:"opiates,omitempty"`
					SixMam          string `json:"sixMam,omitempty"`
					Tramadol        string `json:"tramadol,omitempty"`
				} `json:"results,omitempty"`
			} `json:"drugTestResults,omitempty"`
			PrescribingInformation []string `json:"prescribingInformation,omitempty"`
			RiskInformation        struct {
				LastSelfReportedDate string `json:"lastSelfReportedDate,omitempty"`
				MentalHealthDomain   struct {
					AttemptedSuicide                            string `json:"attemptedSuicide,omitempty"`
					CurrentOrPreviousSelfHarm                   string `json:"currentOrPreviousSelfHarm,omitempty"`
					DiagnosedMentalHealthCondition              string `json:"diagnosedMentalHealthCondition,omitempty"`
					FrequentLifeThreateningSelfHarm             string `json:"frequentLifeThreateningSelfHarm,omitempty"`
					Hallucinations                              string `json:"hallucinations,omitempty"`
					HospitalAdmissionsForMentalHealth           string `json:"hospitalAdmissionsForMentalHealth,omitempty"`
					NoIdentifiedRisk                            string `json:"noIdentifiedRisk,omitempty"`
					NotEngagingWithSupport                      string `json:"notEngagingWithSupport,omitempty"`
					NotTakingPrescribedMedicationAsInstructed   string `json:"notTakingPrescribedMedicationAsInstructed,omitempty"`
					PsychiatricOrPreviousCrisisTeamIntervention string `json:"psychiatricOrPreviousCrisisTeamIntervention,omitempty"`
					Psychosis                                   string `json:"psychosis,omitempty"`
					SelfReportedMentalHealthConcerns            string `json:"selfReportedMentalHealthConcerns,omitempty"`
					ThoughtsOfSuicideOrSelfHarm                 string `json:"thoughtsOfSuicideOrSelfHarm,omitempty"`
				} `json:"mentalHealthDomain,omitempty"`
				RiskOfHarmToSelfDomain struct {
					AssessedAsNotHavingMentalCapacity  string `json:"assessedAsNotHavingMentalCapacity,omitempty"`
					BeliefTheyAreWorthless             string `json:"beliefTheyAreWorthless,omitempty"`
					Hoarding                           string `json:"hoarding,omitempty"`
					LearningDisability                 string `json:"learningDisability,omitempty"`
					MeetsSafeguardingAdultsThreshold   string `json:"meetsSafeguardingAdultsThreshold,omitempty"`
					NoIdentifiedRisk                   string `json:"noIdentifiedRisk,omitempty"`
					OngoingConcernsRelatingToOwnSafety string `json:"ongoingConcernsRelatingToOwnSafety,omitempty"`
					ProblemsMaintainingPersonalHygiene string `json:"problemsMaintainingPersonalHygiene,omitempty"`
					ProblemsMeetingNutritionalNeeds    string `json:"problemsMeetingNutritionalNeeds,omitempty"`
					RequiresIndependentAdvocacy        string `json:"requiresIndependentAdvocacy,omitempty"`
					SelfNeglect                        string `json:"selfNeglect,omitempty"`
				} `json:"riskOfHarmToSelfDomain,omitempty"`
				SocialDomain struct {
					FinancialProblems         string `json:"financialProblems,omitempty"`
					HomelessRoughSleepingNFA  string `json:"homelessRoughSleepingNFA,omitempty"`
					HousingAtRisk             string `json:"housingAtRisk,omitempty"`
					NoIdentifiedRisk          string `json:"noIdentifiedRisk,omitempty"`
					SociallyIsolatedNoSupport string `json:"sociallyIsolatedNoSupport,omitempty"`
				} `json:"socialDomain,omitempty"`
				SubstanceMisuseDomain struct {
					ConfusionOrDisorientation string `json:"ConfusionOrDisorientation,omitempty"`
					AdmissionToAandE          string `json:"admissionToAandE,omitempty"`
					BlackoutOrSeizures        string `json:"blackoutOrSeizures,omitempty"`
					ConcurrentUse             string `json:"concurrentUse,omitempty"`
					HigherRiskDrinking        string `json:"higherRiskDrinking,omitempty"`
					InjectedByOthers          string `json:"injectedByOthers,omitempty"`
					Injecting                 string `json:"injecting,omitempty"`
					InjectingInNeckOrGroin    string `json:"injectingInNeckOrGroin,omitempty"`
					NoIdentifiedRisk          string `json:"noIdentifiedRisk,omitempty"`
					PolyDrugUse               string `json:"polyDrugUse,omitempty"`
					PreviousOverDose          string `json:"previousOverDose,omitempty"`
					RecentPrisonRelease       string `json:"recentPrisonRelease,omitempty"`
					ReducedTolerance          string `json:"reducedTolerance,omitempty"`
					SharingWorks              string `json:"sharingWorks,omitempty"`
					Speedballing              string `json:"speedballing,omitempty"`
					UsingOnTop                string `json:"usingOnTop,omitempty"`
				} `json:"substanceMisuseDomain,omitempty"`
			} `json:"riskInformation,omitempty"`
			SafeguardingInformation struct {
				LastReviewDate     string `json:"lastReviewDate,omitempty"`
				RiskHarmFromOthers string `json:"riskHarmFromOthers,omitempty"`
				RiskToAdults       string `json:"riskToAdults,omitempty"`
				RiskToChildrenOrYP string `json:"riskToChildrenOrYP,omitempty"`
				RiskToSelf         string `json:"riskToSelf,omitempty"`
			} `json:"safeguardingInformation,omitempty"`
		} `json:"client,omitempty"`
		KeyWorker struct {
			LocalIdentifier int `json:"localIdentifier,omitempty"`
			Name            struct {
				Family string `json:"family,omitempty"`
				Given  string `json:"given,omitempty"`
			} `json:"name"`
			Telecom string `json:"telecom,omitempty"`
		} `json:"keyWorker,omitempty"`
	} `json:"data,omitempty"`
}
type PDQv3Response struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"S,attr"`
	Env     string   `xml:"env,attr"`
	Header  struct {
		Text   string `xml:",chardata"`
		Action struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
		} `xml:"Action"`
		MessageID struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
		} `xml:"MessageID"`
		RelatesTo struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
		} `xml:"RelatesTo"`
		To struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
		} `xml:"To"`
	} `xml:"Header"`
	Body struct {
		Text             string `xml:",chardata"`
		PRPAIN201306UV02 struct {
			Text       string `xml:",chardata"`
			Xmlns      string `xml:"xmlns,attr"`
			ITSVersion string `xml:"ITSVersion,attr"`
			ID         struct {
				Text      string `xml:",chardata"`
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"id"`
			CreationTime struct {
				Text  string `xml:",chardata"`
				Value string `xml:"value,attr"`
			} `xml:"creationTime"`
			VersionCode struct {
				Text string `xml:",chardata"`
				Code string `xml:"code,attr"`
			} `xml:"versionCode"`
			InteractionId struct {
				Text      string `xml:",chardata"`
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"interactionId"`
			ProcessingCode struct {
				Text string `xml:",chardata"`
				Code string `xml:"code,attr"`
			} `xml:"processingCode"`
			ProcessingModeCode struct {
				Text string `xml:",chardata"`
				Code string `xml:"code,attr"`
			} `xml:"processingModeCode"`
			AcceptAckCode struct {
				Text string `xml:",chardata"`
				Code string `xml:"code,attr"`
			} `xml:"acceptAckCode"`
			Receiver struct {
				Text     string `xml:",chardata"`
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					Text           string `xml:",chardata"`
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						Text                   string `xml:",chardata"`
						AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
						Root                   string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						Text                    string `xml:",chardata"`
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							Text           string `xml:",chardata"`
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								Text                   string `xml:",chardata"`
								AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								Root                   string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"receiver"`
			Sender struct {
				Text     string `xml:",chardata"`
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					Text           string `xml:",chardata"`
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						Text string `xml:",chardata"`
						Root string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						Text                    string `xml:",chardata"`
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							Text           string `xml:",chardata"`
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								Text string `xml:",chardata"`
								Root string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"sender"`
			Acknowledgement struct {
				Text     string `xml:",chardata"`
				TypeCode struct {
					Text string `xml:",chardata"`
					Code string `xml:"code,attr"`
				} `xml:"typeCode"`
				TargetMessage struct {
					Text string `xml:",chardata"`
					ID   struct {
						Text      string `xml:",chardata"`
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"id"`
				} `xml:"targetMessage"`
			} `xml:"acknowledgement"`
			ControlActProcess struct {
				Text      string `xml:",chardata"`
				ClassCode string `xml:"classCode,attr"`
				MoodCode  string `xml:"moodCode,attr"`
				Code      struct {
					Text       string `xml:",chardata"`
					Code       string `xml:"code,attr"`
					CodeSystem string `xml:"codeSystem,attr"`
				} `xml:"code"`
				Subject struct {
					Text                 string `xml:",chardata"`
					ContextConductionInd string `xml:"contextConductionInd,attr"`
					TypeCode             string `xml:"typeCode,attr"`
					RegistrationEvent    struct {
						Text      string `xml:",chardata"`
						ClassCode string `xml:"classCode,attr"`
						MoodCode  string `xml:"moodCode,attr"`
						ID        struct {
							Text       string `xml:",chardata"`
							NullFlavor string `xml:"nullFlavor,attr"`
						} `xml:"id"`
						StatusCode struct {
							Text string `xml:",chardata"`
							Code string `xml:"code,attr"`
						} `xml:"statusCode"`
						Subject1 struct {
							Text     string `xml:",chardata"`
							TypeCode string `xml:"typeCode,attr"`
							Patient  struct {
								Text      string `xml:",chardata"`
								ClassCode string `xml:"classCode,attr"`
								ID        []struct {
									Text                   string `xml:",chardata"`
									AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
									Extension              string `xml:"extension,attr"`
									Root                   string `xml:"root,attr"`
								} `xml:"id"`
								StatusCode struct {
									Text string `xml:",chardata"`
									Code string `xml:"code,attr"`
								} `xml:"statusCode"`
								EffectiveTime struct {
									Text  string `xml:",chardata"`
									Value string `xml:"value,attr"`
								} `xml:"effectiveTime"`
								PatientPerson struct {
									Text           string `xml:",chardata"`
									ClassCode      string `xml:"classCode,attr"`
									DeterminerCode string `xml:"determinerCode,attr"`
									Name           struct {
										Text   string `xml:",chardata"`
										Use    string `xml:"use,attr"`
										Given  string `xml:"given"`
										Family string `xml:"family"`
									} `xml:"name"`
									Telecom []struct {
										Text  string `xml:",chardata"`
										Use   string `xml:"use,attr"`
										Value string `xml:"value,attr"`
									} `xml:"telecom"`
									AdministrativeGenderCode struct {
										Text           string `xml:",chardata"`
										Code           string `xml:"code,attr"`
										CodeSystem     string `xml:"codeSystem,attr"`
										CodeSystemName string `xml:"codeSystemName,attr"`
									} `xml:"administrativeGenderCode"`
									BirthTime struct {
										Text  string `xml:",chardata"`
										Value string `xml:"value,attr"`
									} `xml:"birthTime"`
									DeceasedInd struct {
										Text  string `xml:",chardata"`
										Value string `xml:"value,attr"`
									} `xml:"deceasedInd"`
									MultipleBirthInd struct {
										Text  string `xml:",chardata"`
										Value string `xml:"value,attr"`
									} `xml:"multipleBirthInd"`
									Addr struct {
										Text              string `xml:",chardata"`
										StreetAddressLine string `xml:"streetAddressLine"`
										City              string `xml:"city"`
										State             string `xml:"state"`
										PostalCode        string `xml:"postalCode"`
										Country           string `xml:"country"`
									} `xml:"addr"`
									MaritalStatusCode struct {
										Text           string `xml:",chardata"`
										Code           string `xml:"code,attr"`
										CodeSystem     string `xml:"codeSystem,attr"`
										CodeSystemName string `xml:"codeSystemName,attr"`
									} `xml:"maritalStatusCode"`
									AsCitizen struct {
										Text            string `xml:",chardata"`
										ClassCode       string `xml:"classCode,attr"`
										PoliticalNation struct {
											Text           string `xml:",chardata"`
											ClassCode      string `xml:"classCode,attr"`
											DeterminerCode string `xml:"determinerCode,attr"`
											Code           struct {
												Text string `xml:",chardata"`
												Code string `xml:"code,attr"`
											} `xml:"code"`
										} `xml:"politicalNation"`
									} `xml:"asCitizen"`
									AsMember struct {
										Text      string `xml:",chardata"`
										ClassCode string `xml:"classCode,attr"`
										Group     struct {
											Text           string `xml:",chardata"`
											ClassCode      string `xml:"classCode,attr"`
											DeterminerCode string `xml:"determinerCode,attr"`
											Code           struct {
												Text           string `xml:",chardata"`
												Code           string `xml:"code,attr"`
												CodeSystem     string `xml:"codeSystem,attr"`
												CodeSystemName string `xml:"codeSystemName,attr"`
											} `xml:"code"`
										} `xml:"group"`
									} `xml:"asMember"`
									BirthPlace struct {
										Text string `xml:",chardata"`
										Addr struct {
											Text string `xml:",chardata"`
											City string `xml:"city"`
										} `xml:"addr"`
									} `xml:"birthPlace"`
								} `xml:"patientPerson"`
								ProviderOrganization struct {
									Text           string `xml:",chardata"`
									ClassCode      string `xml:"classCode,attr"`
									DeterminerCode string `xml:"determinerCode,attr"`
									ID             struct {
										Text       string `xml:",chardata"`
										NullFlavor string `xml:"nullFlavor,attr"`
									} `xml:"id"`
									ContactParty struct {
										Text      string `xml:",chardata"`
										ClassCode string `xml:"classCode,attr"`
									} `xml:"contactParty"`
								} `xml:"providerOrganization"`
								SubjectOf1 struct {
									Text                  string `xml:",chardata"`
									TypeCode              string `xml:"typeCode,attr"`
									QueryMatchObservation struct {
										Text      string `xml:",chardata"`
										ClassCode string `xml:"classCode,attr"`
										MoodCode  string `xml:"moodCode,attr"`
										Code      struct {
											Text       string `xml:",chardata"`
											Code       string `xml:"code,attr"`
											CodeSystem string `xml:"codeSystem,attr"`
										} `xml:"code"`
										Value struct {
											Text  string `xml:",chardata"`
											Xsi   string `xml:"xsi,attr"`
											Value string `xml:"value,attr"`
											Type  string `xml:"type,attr"`
										} `xml:"value"`
									} `xml:"queryMatchObservation"`
								} `xml:"subjectOf1"`
							} `xml:"patient"`
						} `xml:"subject1"`
						Custodian struct {
							Text           string `xml:",chardata"`
							TypeCode       string `xml:"typeCode,attr"`
							AssignedEntity struct {
								Text      string `xml:",chardata"`
								ClassCode string `xml:"classCode,attr"`
								ID        struct {
									Text       string `xml:",chardata"`
									NullFlavor string `xml:"nullFlavor,attr"`
								} `xml:"id"`
							} `xml:"assignedEntity"`
						} `xml:"custodian"`
					} `xml:"registrationEvent"`
				} `xml:"subject"`
				QueryAck struct {
					Text    string `xml:",chardata"`
					QueryId struct {
						Text      string `xml:",chardata"`
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Text string `xml:",chardata"`
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					QueryResponseCode struct {
						Text string `xml:",chardata"`
						Code string `xml:"code,attr"`
					} `xml:"queryResponseCode"`
					ResultTotalQuantity struct {
						Text  string `xml:",chardata"`
						Value string `xml:"value,attr"`
					} `xml:"resultTotalQuantity"`
					ResultCurrentQuantity struct {
						Text  string `xml:",chardata"`
						Value string `xml:"value,attr"`
					} `xml:"resultCurrentQuantity"`
					ResultRemainingQuantity struct {
						Text  string `xml:",chardata"`
						Value string `xml:"value,attr"`
					} `xml:"resultRemainingQuantity"`
				} `xml:"queryAck"`
				QueryByParameter struct {
					Text    string `xml:",chardata"`
					QueryId struct {
						Text      string `xml:",chardata"`
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Text string `xml:",chardata"`
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					ResponseModalityCode struct {
						Text string `xml:",chardata"`
						Code string `xml:"code,attr"`
					} `xml:"responseModalityCode"`
					ResponsePriorityCode struct {
						Text string `xml:",chardata"`
						Code string `xml:"code,attr"`
					} `xml:"responsePriorityCode"`
					MatchCriterionList string `xml:"matchCriterionList"`
					ParameterList      struct {
						Text            string `xml:",chardata"`
						LivingSubjectId struct {
							Text  string `xml:",chardata"`
							Value struct {
								Text      string `xml:",chardata"`
								Extension string `xml:"extension,attr"`
							} `xml:"value"`
							SemanticsText string `xml:"semanticsText"`
						} `xml:"livingSubjectId"`
					} `xml:"parameterList"`
				} `xml:"queryByParameter"`
			} `xml:"controlActProcess"`
		} `xml:"PRPA_IN201306UV02"`
	} `xml:"Body"`
}
type PIXv3Response struct {
	XMLName xml.Name `xml:"Envelope"`
	S       string   `xml:"S,attr"`
	Env     string   `xml:"env,attr"`
	Header  struct {
		Action struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"Action"`
		MessageID struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"MessageID"`
		RelatesTo struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"RelatesTo"`
		To struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"To"`
	} `xml:"Header"`
	Body struct {
		PRPAIN201310UV02 struct {
			Xmlns      string `xml:"xmlns,attr"`
			ITSVersion string `xml:"ITSVersion,attr"`
			ID         struct {
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"id"`
			CreationTime struct {
				Value string `xml:"value,attr"`
			} `xml:"creationTime"`
			VersionCode struct {
				Code string `xml:"code,attr"`
			} `xml:"versionCode"`
			InteractionId struct {
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"interactionId"`
			ProcessingCode struct {
				Code string `xml:"code,attr"`
			} `xml:"processingCode"`
			ProcessingModeCode struct {
				Code string `xml:"code,attr"`
			} `xml:"processingModeCode"`
			AcceptAckCode struct {
				Code string `xml:"code,attr"`
			} `xml:"acceptAckCode"`
			Receiver struct {
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
						Root                   string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								Root                   string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"receiver"`
			Sender struct {
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						Root string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								Root string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"sender"`
			Acknowledgement struct {
				TypeCode struct {
					Code string `xml:"code,attr"`
				} `xml:"typeCode"`
				TargetMessage struct {
					ID struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"id"`
				} `xml:"targetMessage"`
			} `xml:"acknowledgement"`
			ControlActProcess struct {
				ClassCode string `xml:"classCode,attr"`
				MoodCode  string `xml:"moodCode,attr"`
				Code      struct {
					Code       string `xml:"code,attr"`
					CodeSystem string `xml:"codeSystem,attr"`
				} `xml:"code"`
				Subject struct {
					TypeCode          string `xml:"typeCode,attr"`
					RegistrationEvent struct {
						ClassCode string `xml:"classCode,attr"`
						MoodCode  string `xml:"moodCode,attr"`
						ID        struct {
							NullFlavor string `xml:"nullFlavor,attr"`
						} `xml:"id"`
						StatusCode struct {
							Code string `xml:"code,attr"`
						} `xml:"statusCode"`
						Subject1 struct {
							TypeCode string `xml:"typeCode,attr"`
							Patient  struct {
								ClassCode string `xml:"classCode,attr"`
								ID        []struct {
									Extension              string `xml:"extension,attr"`
									Root                   string `xml:"root,attr"`
									AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								} `xml:"id"`
								StatusCode struct {
									Code string `xml:"code,attr"`
								} `xml:"statusCode"`
								PatientPerson struct {
									ClassCode      string `xml:"classCode,attr"`
									DeterminerCode string `xml:"determinerCode,attr"`
									Name           struct {
										Given  string `xml:"given"`
										Family string `xml:"family"`
									} `xml:"name"`
								} `xml:"patientPerson"`
							} `xml:"patient"`
						} `xml:"subject1"`
						Custodian struct {
							TypeCode       string `xml:"typeCode,attr"`
							AssignedEntity struct {
								ClassCode string `xml:"classCode,attr"`
								ID        struct {
									Root string `xml:"root,attr"`
								} `xml:"id"`
							} `xml:"assignedEntity"`
						} `xml:"custodian"`
					} `xml:"registrationEvent"`
				} `xml:"subject"`
				QueryAck struct {
					QueryId struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					QueryResponseCode struct {
						Code string `xml:"code,attr"`
					} `xml:"queryResponseCode"`
					ResultTotalQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultTotalQuantity"`
					ResultCurrentQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultCurrentQuantity"`
					ResultRemainingQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultRemainingQuantity"`
				} `xml:"queryAck"`
				QueryByParameter struct {
					QueryId struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					ResponsePriorityCode struct {
						Code string `xml:"code,attr"`
					} `xml:"responsePriorityCode"`
					ParameterList struct {
						PatientIdentifier struct {
							Value struct {
								AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								Extension              string `xml:"extension,attr"`
								Root                   string `xml:"root,attr"`
							} `xml:"value"`
							SemanticsText string `xml:"semanticsText"`
						} `xml:"patientIdentifier"`
					} `xml:"parameterList"`
				} `xml:"queryByParameter"`
			} `xml:"controlActProcess"`
		} `xml:"PRPA_IN201310UV02"`
	} `xml:"Body"`
}
type PIXmResponse struct {
	ResourceType string `json:"resourceType"`
	ID           string `json:"id"`
	Type         string `json:"type"`
	Total        int    `json:"total"`
	Link         []struct {
		Relation string `json:"relation"`
		URL      string `json:"url"`
	} `json:"link"`
	Entry []struct {
		FullURL  string `json:"fullUrl"`
		Resource struct {
			ResourceType string `json:"resourceType"`
			ID           string `json:"id"`
			Identifier   []struct {
				Use    string `json:"use,omitempty"`
				System string `json:"system"`
				Value  string `json:"value"`
			} `json:"identifier"`
			Active bool `json:"active"`
			Name   []struct {
				Use    string   `json:"use"`
				Family string   `json:"family"`
				Given  []string `json:"given"`
			} `json:"name"`
			Gender    string `json:"gender"`
			BirthDate string `json:"birthDate"`
			Address   []struct {
				Use        string   `json:"use"`
				Line       []string `json:"line"`
				City       string   `json:"city"`
				PostalCode string   `json:"postalCode"`
				Country    string   `json:"country"`
			} `json:"address"`
		} `json:"resource"`
	} `json:"entry"`
}
type TUKPatient struct {
	PIDOID     string `json:"pidoid"`
	PID        string `json:"pid"`
	REGOID     string `json:"regoid"`
	REGID      string `json:"regid"`
	NHSOID     string `json:"nhsoid"`
	NHSID      string `json:"nhsid"`
	GivenName  string `json:"givenname"`
	FamilyName string `json:"familyname"`
	Gender     string `json:"gender"`
	BirthDate  string `json:"birthdate"`
	Street     string `json:"street"`
	Town       string `json:"town"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Zip        string `json:"zip"`
}
type PDQInterface interface {
	pdq() error
}

var (
	pat_cache = make(map[string][]byte)
	DebugMode = false
)

func New_Transaction(i PDQInterface) error {
	return i.pdq()
}
func (i *PDQQuery) pdq() error {
	if err := i.setPDQ_ID(); err != nil {
		return err
	}
	return i.setPatient()
}
func (i *PDQQuery) setPDQ_ID() error {
	if i.Server_URL == "" {
		return errors.New("invalid request - pdq server url is not set")
	}
	if i.REG_OID == "" {
		if os.Getenv(tukcnst.ENV_REG_OID) == "" {
			return errors.New("invalid request - reg oid is not set")
		}
	}
	if i.Timeout == 0 {
		i.Timeout = 5
	}
	if i.NHS_OID == "" {
		i.NHS_OID = tukcnst.NHS_OID_DEFAULT
	}
	if i.MRN_ID != "" && i.MRN_OID != "" {
		i.Used_PID = i.MRN_ID
		i.Used_PID_OID = i.MRN_OID
	} else {
		if i.NHS_ID != "" {
			i.Used_PID = i.NHS_ID
			i.Used_PID_OID = i.NHS_OID
		} else {
			if i.REG_ID != "" && i.REG_OID != "" {
				i.Used_PID = i.REG_ID
				i.Used_PID_OID = i.REG_OID
			}
		}
	}
	if i.Used_PID == "" || i.Used_PID_OID == "" {
		return errors.New("invalid request - no suitable id and oid input values found which can be used for pdq query")
	}
	return nil
}
func (i *PDQQuery) setPatient() error {
	if i.Cache && i.Server_Mode != tukcnst.PDQ_SERVER_TYPE_CGL {
		if _, ok := pat_cache[i.Used_PID]; ok {
			log.Printf("Cache entry found for Patient ID %s", i.Used_PID)
			i.StatusCode = http.StatusOK
			i.Response = pat_cache[i.Used_PID]
			return nil
		}
	}
	var tmplt *template.Template
	var err error
	i.StatusCode = http.StatusOK
	switch i.Server_Mode {
	case tukcnst.PDQ_SERVER_TYPE_CGL:
		i.Request = []byte(i.Server_URL + i.NHS_ID)
		httpReq := tukhttp.CGLRequest{
			Request:      i.Server_URL + i.NHS_ID,
			X_Api_Key:    i.CGL_X_Api_Key,
			X_Api_Secret: i.CGL_X_Api_Secret,
		}
		if err = tukhttp.NewRequest(&httpReq); err == nil {
			if httpReq.StatusCode == http.StatusOK {
				json.Unmarshal(httpReq.Response, &i.CGLUserResponse)
				i.Count = 1
			}
		}
		i.Response = httpReq.Response
		i.StatusCode = httpReq.StatusCode
	case tukcnst.PDQ_SERVER_TYPE_IHE_PIXV3:
		if tmplt, err = template.New(tukcnst.PDQ_SERVER_TYPE_IHE_PIXV3).Funcs(tukutil.TemplateFuncMap()).Parse(tukcnst.GO_Template_PIX_V3_Request); err == nil {
			var b bytes.Buffer
			if err = tmplt.Execute(&b, i); err == nil {
				i.Request = b.Bytes()
				if err = i.newIHESOAPRequest(tukcnst.SOAP_ACTION_PIXV3_Request); err == nil {
					if err = xml.Unmarshal(i.Response, &i.PIXv3Response); err == nil {
						if i.PIXv3Response.Body.PRPAIN201310UV02.Acknowledgement.TypeCode.Code != "AA" {
							err = errors.New("acknowledgement code not equal aa, received " + i.PIXv3Response.Body.PRPAIN201310UV02.Acknowledgement.TypeCode.Code)
						} else {
							i.Count, _ = strconv.Atoi(i.PIXv3Response.Body.PRPAIN201310UV02.ControlActProcess.QueryAck.ResultTotalQuantity.Value)
							if i.Count > 0 {
								pat := TUKPatient{
									PIDOID: i.MRN_OID,
									PID:    i.MRN_ID,
									REGOID: i.REG_OID,
									REGID:  i.REG_ID,
									NHSOID: i.NHS_OID,
									NHSID:  i.NHS_ID,
								}
								pat.GivenName = i.PIXv3Response.Body.PRPAIN201310UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Name.Given
								pat.FamilyName = i.PIXv3Response.Body.PRPAIN201310UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Name.Family
								for _, pid := range i.PIXv3Response.Body.PRPAIN201310UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.ID {
									switch pid.Root {
									case i.REG_OID:
										pat.REGID = pid.Extension
									case i.NHS_OID:
										pat.NHSID = pid.Extension
									case i.MRN_OID:
										pat.PID = pid.Extension
										pat.PIDOID = i.MRN_OID
									}
								}
								if i.Cache {
									pat_cache[i.Used_PID] = i.Response
								}
							}
						}
					}
				}
			}
		}
	case tukcnst.PDQ_SERVER_TYPE_IHE_PDQV3:
		if tmplt, err = template.New(tukcnst.PDQ_SERVER_TYPE_IHE_PDQV3).Funcs(tukutil.TemplateFuncMap()).Parse(tukcnst.GO_Template_PDQ_V3_Request); err == nil {
			var b bytes.Buffer
			if err = tmplt.Execute(&b, i); err == nil {
				i.Request = b.Bytes()
				if err = i.newIHESOAPRequest(tukcnst.SOAP_ACTION_PDQV3_Request); err == nil {
					if err = xml.Unmarshal(i.Response, &i.PDQv3Response); err == nil {
						if i.PDQv3Response.Body.PRPAIN201306UV02.Acknowledgement.TypeCode.Code != "AA" {
							err = errors.New("acknowledgement code not equal aa, received " + i.PDQv3Response.Body.PRPAIN201306UV02.Acknowledgement.TypeCode.Code)
						} else {
							i.Count, _ = strconv.Atoi(i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.QueryAck.ResultTotalQuantity.Value)
							if i.Count > 0 {
								for _, pid := range i.PDQv3Response.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.ID {
									switch pid.Root {
									case i.REG_OID:
										i.REG_ID = pid.Extension
									case i.NHS_OID:
										i.NHS_ID = pid.Extension
									case i.MRN_OID:
										i.MRN_ID = pid.Extension
									}
								}
								if i.Cache {
									pat_cache[i.Used_PID] = i.Response
								}
							}
						}
					}
				}
			}
		}
	case tukcnst.PDQ_SERVER_TYPE_IHE_PIXM:
		i.Request = []byte(i.Server_URL)
		httpReq := tukhttp.PIXmRequest{
			URL:     i.Server_URL,
			PID_OID: i.Used_PID_OID,
			PID:     i.Used_PID,
			Timeout: i.Timeout,
		}
		err = tukhttp.NewRequest(&httpReq)
		i.Response = httpReq.Response
		i.StatusCode = httpReq.StatusCode
		if err == nil {
			if strings.Contains(string(i.Response), "Error") {
				err = errors.New(string(i.Response))
			} else {
				if err := json.Unmarshal(i.Response, &i.PIXmResponse); err == nil {
					log.Printf("%v Patient Entries in Response", i.PIXmResponse.Total)
					i.Count = i.PIXmResponse.Total
					if i.Count > 0 {
						for cnt := 0; cnt < len(i.PIXmResponse.Entry); cnt++ {
							rsppat := i.PIXmResponse.Entry[cnt]
							for _, id := range rsppat.Resource.Identifier {
								if id.System == tukcnst.URN_OID_PREFIX+i.REG_OID {
									i.REG_ID = id.Value
									log.Printf("Set Reg ID %s %s", i.REG_ID, i.REG_OID)
								}
								if id.Use == "usual" {
									i.MRN_ID = id.Value
									i.MRN_OID = strings.Split(id.System, ":")[2]
									log.Printf("Set PID %s %s", i.MRN_ID, i.MRN_OID)
								}
								if id.System == tukcnst.URN_OID_PREFIX+i.NHS_OID {
									i.NHS_ID = id.Value
									log.Printf("Set NHS ID %s %s", i.NHS_ID, i.NHS_OID)
								}
							}
							gn := ""
							for _, name := range rsppat.Resource.Name {
								for _, n := range name.Given {
									gn = gn + n + " "
								}
							}
							i.GivenName = strings.TrimSuffix(gn, " ")
							i.FamilyName = rsppat.Resource.Name[0].Family
							i.BirthDate = strings.ReplaceAll(rsppat.Resource.BirthDate, "-", "")
							i.Gender = rsppat.Resource.Gender

							if len(rsppat.Resource.Address) > 0 {
								i.Zip = rsppat.Resource.Address[0].PostalCode
								if len(rsppat.Resource.Address[0].Line) > 0 {
									i.Street = rsppat.Resource.Address[0].Line[0]
									if len(rsppat.Resource.Address[0].Line) > 1 {
										i.Town = rsppat.Resource.Address[0].Line[1]
									}
								}
								i.City = rsppat.Resource.Address[0].City
								i.Country = rsppat.Resource.Address[0].Country
							}
						}
						if i.Cache {
							pat_cache[i.Used_PID] = i.Response
						}
					}
				}
			}
		}
	}
	if err != nil {
		log.Println(err.Error())
	}
	return err
}
func (i *PDQQuery) newIHESOAPRequest(soapaction string) error {
	httpReq := tukhttp.SOAPRequest{
		URL:        i.Server_URL,
		SOAPAction: soapaction,
		Body:       i.Request,
		Timeout:    i.Timeout,
	}
	err := tukhttp.NewRequest(&httpReq)
	i.Response = httpReq.Response
	i.StatusCode = httpReq.StatusCode
	return err
}
