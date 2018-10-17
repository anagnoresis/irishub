package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/modules/record"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	recordClient "github.com/irisnet/irishub/client/record"
)

type RecordMetadata struct {
	OwnerAddress sdk.AccAddress
	SubmitTime   int64
	DataHash     string
	DataSize     int64
	//PinedNode    string
}

func GetCmdQureyRecord(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query [record ID]",
		Short: "query specified file with record ID",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			recordID := viper.GetString(flagRecordID)

			res, err := cliCtx.QueryStore([]byte(recordID), storeName)
			if len(res) == 0 || err != nil {
				return fmt.Errorf("Record ID [%s] is not existed", recordID)
			}

			var submitFile record.MsgSubmitFile
			cdc.MustUnmarshalBinary(res, &submitFile)

			recordResponse, err := recordClient.ConvertRecordToRecordOutput(cliCtx, submitFile)
			if err != nil {
				return err
			}

			output, err := wire.MarshalJSONIndent(cdc, recordResponse)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil

		},
	}

	cmd.Flags().String(flagRecordID, "", "record ID for query")

	return cmd
}
