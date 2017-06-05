package v2_test

import (
	"testing"

	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	logtest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

var responseV2 = v2.ResponseDetailsView{
	Status:      200,
	Body:        "body",
	EncodedBody: false,
	Headers: map[string][]string{
		"Test": []string{"headers"},
	},
}

func Test_NewSimulationViewFromResponseBody_CanCreateSimulationFromV2Payload(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {
			"pairs": [
				{
					"response": {
						"status": 200,
						"body": "exact match",
						"encodedBody": false,
						"headers": {
							"Header": [
								"value"
							]
						}
					},
					"request": {
						"destination": {
							"exactMatch": "test-server.com"
						}
					}
				}
			],
			"globalActions": {
				"delays": []
			}
		},
		"meta": {
			"schemaVersion": "v2",
			"hoverflyVersion": "v0.11.0",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`))

	Expect(err).To(BeNil())

	Expect(simulation.RequestResponsePairs).To(HaveLen(1))

	Expect(simulation.RequestResponsePairs[0].Request.Body).To(BeNil())
	Expect(*simulation.RequestResponsePairs[0].Request.Destination.ExactMatch).To(Equal("test-server.com"))
	Expect(simulation.RequestResponsePairs[0].Request.Headers).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].Request.Method).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].Request.Path).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].Request.Query).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].Request.Scheme).To(BeNil())

	Expect(simulation.RequestResponsePairs[0].Response.Body).To(Equal("exact match"))
	Expect(simulation.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulation.RequestResponsePairs[0].Response.Headers).To(HaveKeyWithValue("Header", []string{"value"}))
	Expect(simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))

	Expect(simulation.SchemaVersion).To(Equal("v2"))
	Expect(simulation.HoverflyVersion).To(Equal("v0.11.0"))
	Expect(simulation.TimeExported).To(Equal("2017-02-23T12:43:48Z"))
}

func Test_NewSimulationViewFromResponseBody_WontCreateSimulationIfThereIsNoSchemaVersion(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {},
		"meta": {
			"hoverflyVersion": "v0.11.0",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid JSON, missing \"meta.schemaVersion\" string"))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
}

func Test_NewSimulationViewFromResponseBody_WontBlowUpIfMetaIsMissing(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {}
	}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal(`Invalid JSON, missing "meta" object`))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
}

func Test_NewSimulationViewFromResponseBody_CanCreateSimulationFromV1Payload(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {
			"pairs": [
				{
					"response": {
						"status": 200,
						"body": "exact match",
						"encodedBody": false,
						"headers": {
							"Header": [
								"value"
							]
						}
					},
					"request": {
						"destination":"test-server.com"
					}
				}
			],
			"globalActions": {
				"delays": []
			}
		},
		"meta": {
			"schemaVersion": "v1",
			"hoverflyVersion": "v0.11.0",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`))

	Expect(err).To(BeNil())

	Expect(simulation.RequestResponsePairs).To(HaveLen(1))

	Expect(simulation.RequestResponsePairs[0].Request.Body).To(BeNil())
	Expect(*simulation.RequestResponsePairs[0].Request.Destination.ExactMatch).To(Equal("test-server.com"))
	Expect(simulation.RequestResponsePairs[0].Request.Headers).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].Request.Method).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].Request.Path).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].Request.Query).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].Request.Scheme).To(BeNil())

	Expect(simulation.RequestResponsePairs[0].Response.Body).To(Equal("exact match"))
	Expect(simulation.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulation.RequestResponsePairs[0].Response.Headers).To(HaveKeyWithValue("Header", []string{"value"}))
	Expect(simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))

	Expect(simulation.SchemaVersion).To(Equal("v2"))
	Expect(simulation.HoverflyVersion).To(Equal("v0.11.0"))
	Expect(simulation.TimeExported).To(Equal("2017-02-23T12:43:48Z"))
}

func Test_NewSimulationViewFromResponseBody_WontCreateSimulationFromInvalidV1Simulation(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {
			"pairs": [
				{
					
				}
			]
		},
		"meta": {
			"schemaVersion": "v1",
			"hoverflyVersion": "v0.11.0",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid v1 simulation: request is required, response is required"))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
}

func Test_NewSimulationViewFromResponseBody_WontCreateSimulationFromInvalidJson(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{}{}[^.^]{}{}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid JSON"))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
}

func Test_RequestDetailsViewV1_GetQuery_SortsQueryString(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.RequestDetailsViewV1{
		Query: util.StringToPointer("b=b&a=a"),
	}
	queryString := unit.GetQuery()
	Expect(queryString).ToNot(BeNil())

	Expect(*queryString).To(Equal("a=a&b=b"))
}

func Test_RequestDetailsViewV1_GetQuery_ReturnsNilIfNil(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.RequestDetailsViewV1{
		Query: nil,
	}
	queryString := unit.GetQuery()
	Expect(queryString).To(BeNil())
}

func Test_SimulationViewV1_Upgrade_ReturnsAV2Simulation(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV1{
		v2.DataViewV1{
			RequestResponsePairViewV1: []v2.RequestResponsePairViewV1{
				v2.RequestResponsePairViewV1{
					Request: v2.RequestDetailsViewV1{
						RequestType: util.StringToPointer("recording"),
						Scheme:      util.StringToPointer("http"),
						Body:        util.StringToPointer("body"),
						Destination: util.StringToPointer("destination"),
						Method:      util.StringToPointer("GET"),
						Path:        util.StringToPointer("/path"),
						Query:       util.StringToPointer("query=query"),
						Headers: map[string][]string{
							"Test": []string{"headers"},
						},
					},
					Response: responseV2,
				},
			},
		},
		v2.MetaView{
			SchemaVersion:   "v1",
			HoverflyVersion: "test",
			TimeExported:    "today",
		},
	}

	simulationViewV2 := unit.Upgrade()

	Expect(simulationViewV2.RequestResponsePairs).To(HaveLen(1))

	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Scheme).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("http"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Body).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("body"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Destination).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("destination"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Method).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("GET"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Path).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("/path"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Query).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("query=query"),
	}))
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Headers).To(BeEmpty())

	Expect(simulationViewV2.RequestResponsePairs[0].Response).To(Equal(responseV2))

	Expect(simulationViewV2.SchemaVersion).To(Equal("v2"))
	Expect(simulationViewV2.HoverflyVersion).To(Equal("test"))
	Expect(simulationViewV2.TimeExported).To(Equal("today"))
}

func Test_SimulationViewV1_Upgrade_ReturnsGlobMatchesIfTemplate(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV1{
		v2.DataViewV1{
			RequestResponsePairViewV1: []v2.RequestResponsePairViewV1{
				v2.RequestResponsePairViewV1{
					Request: v2.RequestDetailsViewV1{
						RequestType: util.StringToPointer("template"),
						Scheme:      util.StringToPointer("http"),
						Body:        util.StringToPointer("body"),
						Destination: util.StringToPointer("destination"),
						Method:      util.StringToPointer("GET"),
						Path:        util.StringToPointer("/path"),
						Query:       util.StringToPointer("query=query"),
					},
					Response: responseV2,
				},
			},
		},
		v2.MetaView{
			SchemaVersion:   "v1",
			HoverflyVersion: "test",
			TimeExported:    "today",
		},
	}

	simulationViewV2 := unit.Upgrade()

	Expect(simulationViewV2.RequestResponsePairs).To(HaveLen(1))

	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Scheme).To(Equal(v2.RequestFieldMatchersView{
		GlobMatch: util.StringToPointer("http"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Body).To(Equal(v2.RequestFieldMatchersView{
		GlobMatch: util.StringToPointer("body"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Destination).To(Equal(v2.RequestFieldMatchersView{
		GlobMatch: util.StringToPointer("destination"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Method).To(Equal(v2.RequestFieldMatchersView{
		GlobMatch: util.StringToPointer("GET"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Path).To(Equal(v2.RequestFieldMatchersView{
		GlobMatch: util.StringToPointer("/path"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Query).To(Equal(v2.RequestFieldMatchersView{
		GlobMatch: util.StringToPointer("query=query"),
	}))
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Headers).To(BeEmpty())
}

func Test_SimulationViewV2_Upgrade_UpdatesMetadataToV3(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{}

	simulationViewV3 := unit.Upgrade()
	Expect(simulationViewV3.MetaView.SchemaVersion).To(Equal("v3"))
}

func Test_SimulationViewV2_Upgrade_UpgradesPairs(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		DataViewV2: v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Scheme: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("http"),
						},
						Method: &v2.RequestFieldMatchersView{
							GlobMatch: util.StringToPointer("*"),
						},
						Destination: &v2.RequestFieldMatchersView{
							GlobMatch: util.StringToPointer("*"),
						},
						Path: &v2.RequestFieldMatchersView{
							RegexMatch: util.StringToPointer("api"),
						},

						Body: &v2.RequestFieldMatchersView{
							JsonMatch: util.StringToPointer(`{"api": true}`),
						},
						Headers: map[string][]string{
							"TestHeader": {"one"},
						},
					},
					Response: responseV2,
				},
			},
		},
	}

	simulationViewV3 := unit.Upgrade()
	Expect(simulationViewV3.DataViewV3.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationViewV3.DataViewV3.RequestResponsePairs[0].Request).To(Equal(v2.RequestDetailsViewV3{
		Scheme: &v2.RequestFieldMatchersView{
			ExactMatch: util.StringToPointer("http"),
		},
		Method: &v2.RequestFieldMatchersView{
			GlobMatch: util.StringToPointer("*"),
		},
		Destination: &v2.RequestFieldMatchersView{
			GlobMatch: util.StringToPointer("*"),
		},
		Path: &v2.RequestFieldMatchersView{
			RegexMatch: util.StringToPointer("api"),
		},

		Body: &v2.RequestFieldMatchersView{
			JsonMatch: util.StringToPointer(`{"api": true}`),
		},
		Headers: map[string][]string{
			"TestHeader": {"one"},
		},
	}))
	Expect(simulationViewV3.DataViewV3.RequestResponsePairs[0].Response).To(Equal(responseV2))
}

func Test_SimulationViewV2_Upgrade_UpgradesSingleQueryExactMatch(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		DataViewV2: v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Query: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer(`q=something`),
						},
					},
					Response: responseV2,
				},
			},
		},
	}

	simulationViewV3 := unit.Upgrade()
	Expect(simulationViewV3.DataViewV3.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationViewV3.DataViewV3.RequestResponsePairs[0].Request).To(Equal(v2.RequestDetailsViewV3{
		Query: map[string]*v2.RequestFieldMatchersView{
			"q": &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("something"),
			},
		},
	}))
}

func Test_SimulationViewV2_Upgrade_UpgradesSingleQueryGlobMatch(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		DataViewV2: v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Query: &v2.RequestFieldMatchersView{
							GlobMatch: util.StringToPointer(`q=*`),
						},
					},
					Response: responseV2,
				},
			},
		},
	}

	simulationViewV3 := unit.Upgrade()
	Expect(simulationViewV3.DataViewV3.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationViewV3.DataViewV3.RequestResponsePairs[0].Request).To(Equal(v2.RequestDetailsViewV3{
		Query: map[string]*v2.RequestFieldMatchersView{
			"q": &v2.RequestFieldMatchersView{
				GlobMatch: util.StringToPointer("*"),
			},
		},
	}))
}

func Test_SimulationViewV2_Upgrade_UpgradesMultipleQueryExactMatch(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		DataViewV2: v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Query: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer(`limit=30&order=desc`),
						},
					},
					Response: responseV2,
				},
			},
		},
	}

	simulationViewV3 := unit.Upgrade()
	Expect(simulationViewV3.DataViewV3.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationViewV3.DataViewV3.RequestResponsePairs[0].Request).To(Equal(v2.RequestDetailsViewV3{
		Query: map[string]*v2.RequestFieldMatchersView{
			"limit": &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("30"),
			},
			"order": &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("desc"),
			},
		},
	}))
}

func Test_SimulationViewV2_Upgrade_UpgradesMultipleQueryGlobMatch(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		DataViewV2: v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Query: &v2.RequestFieldMatchersView{
							GlobMatch: util.StringToPointer(`limit=*&order=asc`),
						},
					},
					Response: responseV2,
				},
			},
		},
	}

	simulationViewV3 := unit.Upgrade()
	Expect(simulationViewV3.DataViewV3.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationViewV3.DataViewV3.RequestResponsePairs[0].Request).To(Equal(v2.RequestDetailsViewV3{
		Query: map[string]*v2.RequestFieldMatchersView{
			"limit": &v2.RequestFieldMatchersView{
				GlobMatch: util.StringToPointer("*"),
			},
			"order": &v2.RequestFieldMatchersView{
				GlobMatch: util.StringToPointer("asc"),
			},
		},
	}))
}

func Test_SimulationViewV2_Upgrade_UpgradesSingleQueryExactMatchOnlyKey(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		DataViewV2: v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Query: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer(`something`),
						},
					},
					Response: responseV2,
				},
			},
		},
	}

	simulationViewV3 := unit.Upgrade()

	Expect(simulationViewV3.DataViewV3.RequestResponsePairs[0].Request).To(Equal(v2.RequestDetailsViewV3{
		Query: map[string]*v2.RequestFieldMatchersView{
			"something": &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer(""),
			},
		},
	}))
}

func Test_SimulationViewV2_Upgrade_UpgradesMultipleQueryExactMatchOnlyKey(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		DataViewV2: v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Query: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer(`something&else`),
						},
					},
					Response: responseV2,
				},
			},
		},
	}

	simulationViewV3 := unit.Upgrade()

	Expect(simulationViewV3.DataViewV3.RequestResponsePairs[0].Request).To(Equal(v2.RequestDetailsViewV3{
		Query: map[string]*v2.RequestFieldMatchersView{
			"something": &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer(""),
			},
			"else": &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer(""),
			},
		},
	}))
}

func Test_SimulationViewV2_Upgrade_RespectsCommaSeperatedQueryValue(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		DataViewV2: v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Query: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer(`something=one,two`),
						},
					},
					Response: responseV2,
				},
			},
		},
	}

	simulationViewV3 := unit.Upgrade()

	Expect(simulationViewV3.DataViewV3.RequestResponsePairs[0].Request).To(Equal(v2.RequestDetailsViewV3{
		Query: map[string]*v2.RequestFieldMatchersView{
			"something": &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("one,two"),
			},
		},
	}))
}

func Test_SimulationViewV2_Upgrade_UpgradesMultipleQueryValuesUsesTheLastOne(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV2{
		DataViewV2: v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Query: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer(`something=one&something=two`),
						},
					},
					Response: responseV2,
				},
			},
		},
	}

	simulationViewV3 := unit.Upgrade()

	Expect(simulationViewV3.DataViewV3.RequestResponsePairs[0].Request).To(Equal(v2.RequestDetailsViewV3{
		Query: map[string]*v2.RequestFieldMatchersView{
			"something": &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("two"),
			},
		},
	}))
}

func Test_SimulationViewV2_Upgrade_LogsWhenAnErrorOccursDuringQueryParsing(t *testing.T) {
	RegisterTestingT(t)

	log.SetOutput(ioutil.Discard)
	logger := log.StandardLogger()
	testHook := logtest.NewLocal(logger)

	unit := v2.SimulationViewV2{
		DataViewV2: v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Query: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer(`% % % %`),
						},
					},
					Response: responseV2,
				},
			},
		},
	}

	unit.Upgrade()

	Expect(testHook.Entries).To(HaveLen(1))

	Expect(testHook.Entries[0].Level).To(Equal(log.WarnLevel))
	Expect(testHook.Entries[0].Message).To(Equal("There was an error when upgrading v2 simulation to v3"))
	Expect(testHook.Entries[0].Data["query"]).To(Equal("% % % %"))
}
func Test_SimulationViewV1_Upgrade_CanReturnAnIncompleteRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV1{
		v2.DataViewV1{
			RequestResponsePairViewV1: []v2.RequestResponsePairViewV1{
				v2.RequestResponsePairViewV1{
					Request: v2.RequestDetailsViewV1{
						Method: util.StringToPointer("POST"),
					},
					Response: responseV2,
				},
			},
		},
		v2.MetaView{
			SchemaVersion:   "v1",
			HoverflyVersion: "test",
			TimeExported:    "today",
		},
	}

	simulationViewV2 := unit.Upgrade()

	Expect(simulationViewV2.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationViewV2.RequestResponsePairs[0].Request.Scheme).To(BeNil())
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Body).To(BeNil())
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Destination).To(BeNil())
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Method).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("POST"),
	}))
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Path).To(BeNil())
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Query).To(BeNil())
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Headers).To(BeNil())

	Expect(simulationViewV2.RequestResponsePairs[0].Response.Status).To(Equal(200))
	Expect(simulationViewV2.RequestResponsePairs[0].Response.Body).To(Equal("body"))
	Expect(simulationViewV2.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulationViewV2.RequestResponsePairs[0].Response.Headers).To(HaveKeyWithValue("Test", []string{"headers"}))
}
