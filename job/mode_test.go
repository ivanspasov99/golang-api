package job

//var testGetJobModeWriter = []struct {
//	name         string
//	requestUrl   string
//	expectedFunc ResponseWriter
//}{
//	{
//		"Test with no query mode parameter should use JSON default writer mode",
//		"/job",
//		writeJSON,
//	},
//}
//
//func TestGetJobModeWriter(t *testing.T) {
//	for _, tt := range testGetJobModeWriter {
//		t.Run(tt.name, func(t *testing.T) {
//			r, _ := http.NewRequest("GET", tt.requestUrl, nil)
//
//			w := getJobModeWriter(r)
//
//			if &tt.expectedFunc != &w {
//				t.Error("Received function is not the same as the expected")
//			}
//		})
//	}
//}
