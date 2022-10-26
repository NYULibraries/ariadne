package sfx

// TODO: Redo this test so that it works with new MultiObjectsResponse type, which
// has a JSON field that needs to be tested.  `toResponseJSON` has been deleted.
//func TestToResponseJson(t *testing.T) {
//	dummyGoodXMLResponse := `
//<ctx_obj_set>
//	<ctx_obj>
//		<ctx_obj_targets>
//			<target>
//				<target_url>http://answers.library.newschool.edu/</target_url>
//			</target>
//		</ctx_obj_targets>
//	</ctx_obj>
//</ctx_obj_set>`
//	dummyBadXMLResponse := `
//<ctx_obj_set`
//	dummyJSONResponse := `{
//    "ctx_obj": [
//        {
//            "ctx_obj_targets": [
//                {
//                    "target": [
//                        {
//                            "target_name": "",
//                            "target_public_name": "",
//                            "target_url": "http://answers.library.newschool.edu/",
//                            "authentication": "",
//                            "proxy": ""
//                        }
//                    ]
//                }
//            ]
//        }
//    ]
//}`
//	var tests = []struct {
//		from        []byte
//		expected    string
//		expectedErr error
//	}{
//		{[]byte(dummyGoodXMLResponse), dummyJSONResponse, nil},
//		{[]byte(dummyBadXMLResponse), "", errors.New("error")},
//	}
//
//	for _, tt := range tests {
//		testname := fmt.Sprintf("%s", tt.expected)
//		t.Run(testname, func(t *testing.T) {
//			ans, err := toResponseJSON(tt.from)
//			// if err != nil {
//			// 	t.Errorf("error %v", err)
//			// }
//			if tt.expectedErr != nil {
//				if err == nil {
//					t.Errorf("toResponseJSON err was '%v', expecting '%v'", err, tt.expectedErr)
//				}
//			}
//			if ans != tt.expected {
//				t.Errorf("toResponseJSON was '%v', expecting '%v'", ans, tt.expected)
//			}
//		})
//	}
//}
