package purpleair

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type PurpleAir struct {
	Ids []int64 `toml:"ids"`
}

const sampleConfig = `
  ## List of sensor IDs to monitor.
  ids = [
      5651,
      21887,
      21888,
      4355,
      4356,
      19843,
      19844,
      28299,
      28300
  ]
`

func (p PurpleAir) SampleConfig() string {
	return sampleConfig
}

func (p PurpleAir) Description() string {
	return "Gather information from selected sensors via PurpleAir's JSON API"
}

func (p *PurpleAir) Gather(accumulator telegraf.Accumulator) error {
	fmt.Println("PurpleAir Gatherin'")
	for _, id := range p.Ids {
		url := fmt.Sprintf("https://www.purpleair.com/json?show=%d", id)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer func() { _ = resp.Body.Close() }()

		var sr SensorReading
		dec := json.NewDecoder(resp.Body)
		if err = dec.Decode(&sr); err != nil {
			return err
		}
		fmt.Printf("%+v\n", sr)
	}
	return nil
}

type SensorReading struct {
	Id       int64 `json:"ID"`
	ParentId int64 `json:"ParentID,omitempty"`
	Label    string `json:"Label"`
	Lat      float64 `json:"Lat"`
	Lon      float64 `json:"Lon"`
	Type     string `json:"Type"`
	TempF    float64 `json:"temp_f"`
}

func init() {
	inputs.Add("purpleair", func() telegraf.Input {
		return &PurpleAir{}
	})
}
