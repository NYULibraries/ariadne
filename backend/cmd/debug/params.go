package debug

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
)

func init() {
	DebugCmd.AddCommand(dumpParamsCmd)
}

var dumpParamsCmd = &cobra.Command{
	Use:     "params [query string]",
	Short:   "Dump params in query string as JSON object",
	Example: "ariadne debug params 'url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404<fssessid>0<%2Ffssessid>&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat'",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var queryString = args[0]
		dump, err := dumpParams(queryString)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(dump)
	},
}

func dumpParams(queryString string) (string, error) {
	urlValues, err := url.ParseQuery(queryString)
	if err != nil {
		return queryString, err
	}

	bytes, err := json.MarshalIndent(urlValues, "", "  ")
	if err != nil {
		return queryString, err
	}

	return string(bytes), nil
}
