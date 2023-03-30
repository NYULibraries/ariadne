package debug

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"

	"ariadne/sfx"
)

func init() {
	DebugCmd.AddCommand(dumpSFXHTTPRequestCmd)
	DebugCmd.AddCommand(dumpSFXHTTPResponseCmd)
	DebugCmd.AddCommand(targetsJSONCmd)
}

var dumpSFXHTTPRequestCmd = &cobra.Command{
	Use:     "sfx-request [query string]",
	Short:   "Dump SFX HTTP request for query string",
	Example: "ariadne debug sfx-request 'url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404<fssessid>0<%2Ffssessid>&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat'",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var queryString = args[0]
		dump, err := dumpSFXHTTPRequest(queryString)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(dump)
	},
}

var dumpSFXHTTPResponseCmd = &cobra.Command{
	Use:     "sfx-response [query string]",
	Short:   "Dump SFX HTTP response for query string",
	Example: "ariadne debug sfx-response 'url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404<fssessid>0<%2Ffssessid>&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat'",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var queryString = args[0]
		dump, err := dumpSFXHTTPResponse(queryString)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(dump)
	},
}

func dumpSFXHTTPRequest(queryString string) (string, error) {
	sfxRequest, err := sfx.NewSFXRequest(queryString)
	if err != nil {
		return queryString, err
	}

	return sfxRequest.DumpedHTTPRequest, nil
}

func dumpSFXHTTPResponse(queryString string) (string, error) {
	sfxRequest, err := sfx.NewSFXRequest(queryString)
	if err != nil {
		return queryString, err
	}

	sfxResponse, err := sfx.Do(sfxRequest)
	if err != nil {
		return queryString, err
	}

	return sfxResponse.DumpedHTTPResponse, nil
}

var targetsJSONCmd = &cobra.Command{
	Use:     "sfx-targets [query string]",
	Short:   "Return JSON array of target objects returned by SFX response for query string",
	Example: "ariadne debug sfx-targets 'url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404<fssessid>0<%2Ffssessid>&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat'",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var queryString = args[0]
		targetsJSON, err := targetsJSON(queryString)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(targetsJSON)
	},
}

func targetsJSON(queryString string) (string, error) {
	sfxRequest, err := sfx.NewSFXRequest(queryString)
	if err != nil {
		return queryString, err
	}

	sfxResponse, err := sfx.Do(sfxRequest)
	if err != nil {
		return queryString, err
	}

	targetsJSON, err := json.MarshalIndent(
		(*sfxResponse.XMLResponseBody.ContextObject)[0].SFXContextObjectTargets,
		"",
		"    ")
	if err != nil {
		return "", err
	}

	return string(targetsJSON), nil
}
