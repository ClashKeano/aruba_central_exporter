package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ApResponse struct {
	AccessPoints []AccessPoint `json:"aps"`
}

type AccessPoint struct {
	ApDeploymentNode   string   `json:"ap_deployment_mode"`
	ApGroup            string   `json:"ap_group"`
	ClientCount        int      `json:"client_count"`
	ClusterId          string   `json:"cluster_id"`
	ControllerName     string   `json:"controller_name"`
	CpuUtilization     int      `json:"cpu_utilization"`
	FirmwareVersion    string   `json:"firmware_version"`
	GatewayClusterId   string   `json:"gateway_cluster_id"`
	GatewayClusterName string   `json:"gateway_cluster_name"`
	GroupName          string   `json:"group_name"`
	IpAddress          string   `json:"ip_address"`
	Labels             []string `json:"labels"`
	LastModified       int      `json:"last_modified"`
	MacAddress         string   `json:"macaddr"`
	MemFree            int      `json:"mem_free"`
	MemTotal           int      `json:"mem_total"`
	MeshRole           string   `json:"mesh_role"`
	Model              string   `json:"model"`
	Name               string   `json:"name"`
	Notes              string   `json:"notes"`
	PublicIpAddress    string   `json:"public_ip_address"`
	Radios             []struct {
		Band          int    `json:"band"`
		Channel       string `json:"channel"`
		Index         int    `json:"index"`
		MacAddress    string `json:"macaddr"`
		Node          int    `json:"node"`
		RadioName     string `json:"radio_name"`
		RadioType     string `json:"radio_type"`
		SpatialStream string `json:"spatial_stream"`
		Status        string `json:"status"`
		TxPower       int    `json:"tx_power"`
		Utilization   int    `json:"utilization"`
	} `json:"radios"`
	Serial      string `json:"serial"`
	Site        string `json:"site"`
	SleepStatus bool   `json:"sleep_status"`
	Status      string `json:"status"`
	SubnetMask  string `json:"subnet_mask"`
	SwarmId     string `json:"swarm_id"`
	SwarmMaster bool   `json:"swarm_master"`
	SwarmName   string `json:"swarm_name"`
	Uptime      int    `json:"uptime"`
}

type McResponse struct {
	Count               int                  `json:"count"`
	MobilityControllers []MobilityController `json:"mcs"`
}

type MobilityController struct {
	CpuUtilization        int      `json:"cpu_utilization"`
	FirmwareBackupVersion string   `json:"firmware_backup_version"`
	FirmwareVersion       string   `json:"firmware_version"`
	GroupName             string   `json:"group_name"`
	IpAddress             string   `json:"ip_address"`
	Labels                []string `json:"labels"`
	MacRange              string   `json:"mac_range"`
	MacAddress            string   `json:"macaddr"`
	MemFree               int      `json:"mem_free"`
	MemTotal              int      `json:"mem_total"`
	Mode                  string   `json:"mode"`
	Model                 string   `json:"model"`
	Name                  string   `json:"name"`
	RebootReason          string   `json:"reboot_reason"`
	Role                  string   `json:"role"`
	Serial                string   `json:"serial"`
	Site                  string   `json:"site"`
	Status                string   `json:"status"`
	Uptime                int      `json:"uptime"`
}

type SiteResponse struct {
	Sites []Site `json:"items"`
}

type Site struct {
	BranchCpuHigh          int     `json:"branch_cpu_high"`
	BranchDeviceStatusDown int     `json:"branch_device_status_down"`
	BranchDeviceStatusUp   int     `json:"branch_device_status_up"`
	BranchMemHigh          int     `json:"branch_mem_high"`
	CapeState              string  `json:"cape_state"`
	CapeStateDesc          string  `json:"cape_state_dscr"`
	CapeStateUrl           string  `json:"cape_url"`
	ConnectedCount         int     `json:"connected_count"`
	DeviceDown             int     `json:"device_down"`
	DeviceHighCh24         int     `json:"device_high_ch_2_4ghz"`
	DeviceHighCh5          int     `json:"device_high_ch_5ghz"`
	DeviceHighCpu          int     `json:"device_high_cpu"`
	DeviceHighMem          int     `json:"device_high_mem"`
	DeviceHighNoise24      int     `json:"device_high_noise_2_4ghz"`
	DeviceHighNoise5       int     `json:"device_high_noise_5ghz"`
	DeviceUp               int     `json:"device_up"`
	FailedCount            int     `json:"failed_count"`
	Id                     string  `json:"id"`
	InsightHi              int     `json:"insight_hi"`
	InsightLo              int     `json:"insight_lo"`
	InsightMi              int     `json:"insight_mi"`
	Lat                    float64 `json:"lat"`
	Long                   float64 `json:"long"`
	Name                   string  `json:"name"`
	PotentialIssue         bool    `json:"potential_issue"`
	Score                  int     `json:"score"`
	SilverPeakState        string  `json:"silverpeak_state"`
	SilverPeakStateSummary string  `json:"silverpeak_state_summary"`
	SilverPeakUrl          string  `json:"silverpeak_url"`
	UserConnHealthScore    int     `json:"user_conn_health_score"`
	WanTunnelsDown         int     `json:"wan_tunnels_down"`
	WanTunnelsNoIssue      int     `json:"wan_tunnels_no_issue"`
	WanUplinksDown         int     `json:"wan_uplinks_down"`
	WanUplinksNoIssue      int     `json:"wan_uplinks_no_issue"`
	WiredCPUHigh           int     `json:"wired_cpu_high"`
	WiredDeviceSatusDown   int     `json:"wired_device_status_down"`
	WiredDeviceStatusUp    int     `json:"wired_device_status_up"`
	WiredMemHigh           int     `json:"wired_mem_high"`
	WlanCpuHigh            int     `json:"wlan_cpu_high"`
	WlanDeviceStatusDown   int     `json:"wlan_device_status_down"`
	WlanDeviceStatusUp     int     `json:"wlan_device_status_up"`
	WlanMemHigh            int     `json:"wlan_mem_high"`
}

type SwitchResponse struct {
	Count    int      `json:"count"`
	Switches []Switch `json:"switches"`
}

type Switch struct {
	ClientCount      int      `json:"client_count"`
	CPUUtilization   int      `json:"cpu_utilization"`
	FanSpeed         string   `json:"fan_speed"`
	FirmwareVersion  string   `json:"firmware_version"`
	GroupID          int      `json:"group_id"`
	GroupName        string   `json:"group_name"`
	IPAddress        string   `json:"ip_address"`
	LabelIDs         []int    `json:"label_ids"`
	Labels           []string `json:"labels"`
	MacAddress       string   `json:"macaddr"`
	MaxPower         int      `json:"max_power"`
	MemFree          int      `json:"mem_free"`
	MemTotal         int      `json:"mem_total"`
	Model            string   `json:"model"`
	Name             string   `json:"name"`
	PoeConsumption   string   `json:"poe_consumption"`
	PowerConsumption int      `json:"power_consumption"`
	PublicIPAddress  string   `json:"public_ip_address"`
	Serial           string   `json:"serial"`
	Site             string   `json:"site"`
	SiteID           int      `json:"site_id"`
	StackID          string   `json:"stack_id"`
	StackMemberID    int      `json:"stack_member_id"`
	Status           string   `json:"status"`
	SwitchRole       int      `json:"switch_role"`
	SwitchType       string   `json:"switch_type"`
	Temperature      string   `json:"temperature"`
	UplinkPorts      []struct {
		Port string `json:"port"`
	} `json:"uplink_ports"`
	Uptime int `json:"uptime"`
	Usage  int `json:"usage"`
}

type TokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type TopNClientResponse struct {
	Clients []Client `json:"clients"`
}

type Client struct {
	MacAddress  string `json:"macaddr"`
	Name        string `json:"name"`
	RxDataBytes int    `json:"rx_data_bytes"`
	TxDataBytes int    `json:"tx_data_bytes"`
}

var (
	apClientCount    = prometheus.NewDesc("aruba_ap_client_count", "Number of clients connected to access point", []string{"name", "groupName", "site", "status", "firmwareVersion", "model"}, nil)
	apCpuUtilization = prometheus.NewDesc("aruba_ap_cpu_utilization", "CPU Utilization of the access point in percentge", []string{"name", "groupName", "site", "status", "firmwareVersion", "model"}, nil)
	apMemFree        = prometheus.NewDesc("aruba_ap_mem_free", "Amount of free memory of access point", []string{"name", "groupName", "site", "status", "firmwareVersion", "model"}, nil)
	apMemTotal       = prometheus.NewDesc("aruba_ap_mem_total", "Total amount of  memory of access point", []string{"name", "groupName", "site", "status", "firmwareVersion", "model"}, nil)
	apUptime         = prometheus.NewDesc("aruba_ap_uptime", "Uptime of the access point in seconds", []string{"name", "groupName", "site", "status", "firmwareVersion", "model"}, nil)

	apRadioTxPower     = prometheus.NewDesc("aruba_ap_radio_tx_power", "Radio tx power", []string{"band", "channel", "radioName", "apName"}, nil)
	apRadioUtilization = prometheus.NewDesc("aruba_ap_radio_utilization", "Radip cpu utilization", []string{"band", "channel", "radioName", "apName"}, nil)

	clientRxDataBytes = prometheus.NewDesc("aruba_client_rx_data_bytes", "Volume of data received", []string{"name", "mac"}, nil)
	clientTxDataBytes = prometheus.NewDesc("aruba_client_tx_data_bytes", "Volume of data transmitted", []string{"name", "mac"}, nil)

	mcCpuUtilization = prometheus.NewDesc("aruba_mc_cpu_utilization", "CPU Utilization of the mobility controller in percentge", []string{"name", "groupName", "mode", "model", "site", "status", "firmwareVersion"}, nil)
	mcMemFree        = prometheus.NewDesc("aruba_mc_mem_free", "Amount of free memory of mobility controller", []string{"name", "groupName", "mode", "model", "site", "status", "firmwareVersion"}, nil)
	mcMemTotal       = prometheus.NewDesc("aruba_mc_mem_total", "Total amount of  memory of mobility controller", []string{"name", "groupName", "mode", "model", "site", "status", "firmwareVersion"}, nil)
	mcUptime         = prometheus.NewDesc("aruba_mc_uptime", "Uptime of the mobility controller in seconds", []string{"name", "groupName", "mode", "model", "site", "status", "firmwareVersion"}, nil)

	siteConnectedCount        = prometheus.NewDesc("aruba_site_connected_count", "Number of connected devices", []string{"name", "id"}, nil)
	siteDeviceDown            = prometheus.NewDesc("aruba_site_device_down", "Number of down devices", []string{"name", "id"}, nil)
	siteDeviceHighCh24        = prometheus.NewDesc("aruba_site_device_high_ch_2_4ghz", "Number of devices with high 2.4ghz channel utilization", []string{"name", "id"}, nil)
	siteDeviceHighCh5         = prometheus.NewDesc("aruba_site_device_high_ch_5ghz", "Number of devices with high 5ghz channel utilization", []string{"name", "id"}, nil)
	siteDeviceHighCpu         = prometheus.NewDesc("aruba_site_device_high_cpu", "Number of devices with high cpu utilization", []string{"name", "id"}, nil)
	siteDeviceHighMem         = prometheus.NewDesc("aruba_site_device_high_mem", "Number of devices with high mem utilization", []string{"name", "id"}, nil)
	siteDeviceHighNoise24     = prometheus.NewDesc("aruba_site_device_high_noise_2_4ghz", "Number of devices with high 2.4ghz noise", []string{"name", "id"}, nil)
	siteDeviceHighNoise5      = prometheus.NewDesc("aruba_site_device_high_noise_5ghz", "Number of devices with high 5ghz noise", []string{"name", "id"}, nil)
	siteDeviceUp              = prometheus.NewDesc("aruba_site_device_up", "Number of up devices", []string{"name", "id"}, nil)
	siteWiredCpuHigh          = prometheus.NewDesc("aruba_site_wired_cpu_high", "Number of wired devices with high CPU", []string{"name", "id"}, nil)
	siteWiredDeviceStatusDown = prometheus.NewDesc("aruba_site_wired_device_status_down", "Number of wired devices up", []string{"name", "id"}, nil)
	siteWiredDeviceStatusUp   = prometheus.NewDesc("aruba_site_wired_device_status_up", "Number of wired devices down", []string{"name", "id"}, nil)
	siteWiredMemHigh          = prometheus.NewDesc("aruba_site_wired_mem_high", "Number of wired devices with high memory usage", []string{"name", "id"}, nil)
	siteWlanCpuHigh           = prometheus.NewDesc("aruba_site_wlan_cpu_high", "Number of wired devices with high cpu usage", []string{"name", "id"}, nil)
	siteWlanDeviceStatusDown  = prometheus.NewDesc("aruba_site_wlan_device_status_down", "Number of down wireless devices", []string{"name", "id"}, nil)
	siteWlanDeviceStatusUp    = prometheus.NewDesc("aruba_site_wlan_device_status_up", "Number of down wireless devices", []string{"name", "id"}, nil)
	siteWlanMemHigh           = prometheus.NewDesc("aruba_site_wlan_mem_high", "Number of wireless devices with high cpu", []string{"name", "id"}, nil)

	switchClientCount    = prometheus.NewDesc("aruba_switch_client_count", "Number of clients connected to switch", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion", "model"}, nil)
	switchCpuUtilization = prometheus.NewDesc("aruba_switch_cpu_utilization", "Current Switch CPU utilization percentage", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion", "model"}, nil)
	switchMemFree        = prometheus.NewDesc("aruba_switch_mem_free", "Switch free memory", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion", "model"}, nil)
	switchMemTotal       = prometheus.NewDesc("aruba_switch_mem_total", "Switch total memory", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion", "model"}, nil)
	switchUsage          = prometheus.NewDesc("aruba_switch_usage", "Switch uptime", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion", "model"}, nil)
	switchUptime         = prometheus.NewDesc("aruba_switch_uptime", "Switch usage", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion", "model"}, nil)

	expiresIn  = 0
	configFile string
	verbose    bool
	version    = 1.1
)

type Exporter struct {
	arubaEndpoint, arubaAccessToken, arubaRefreshToken string
}

func NewExporter(arubaEndpoint string, arubaAccessToken string, arubaRefreshToken string) *Exporter {
	return &Exporter{
		arubaEndpoint:     arubaEndpoint,
		arubaAccessToken:  arubaAccessToken,
		arubaRefreshToken: arubaRefreshToken,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- apClientCount
	ch <- apCpuUtilization
	ch <- apMemFree
	ch <- apMemTotal
	ch <- apUptime

	ch <- apRadioTxPower
	ch <- apRadioUtilization

	ch <- clientRxDataBytes
	ch <- clientTxDataBytes

	ch <- mcCpuUtilization
	ch <- mcMemFree
	ch <- mcMemTotal
	ch <- mcUptime

	ch <- siteConnectedCount
	ch <- siteDeviceDown
	ch <- siteDeviceHighCh24
	ch <- siteDeviceHighCh5
	ch <- siteDeviceHighCpu
	ch <- siteDeviceHighMem
	ch <- siteDeviceHighNoise24
	ch <- siteDeviceHighNoise5
	ch <- siteDeviceUp
	ch <- siteWiredCpuHigh
	ch <- siteWiredDeviceStatusDown
	ch <- siteWiredDeviceStatusUp
	ch <- siteWiredMemHigh
	ch <- siteWlanCpuHigh
	ch <- siteWlanDeviceStatusDown
	ch <- siteWlanDeviceStatusUp
	ch <- siteWlanMemHigh

	ch <- switchClientCount
	ch <- switchCpuUtilization
	ch <- switchMemFree
	ch <- switchMemTotal
	ch <- switchUsage
	ch <- switchUptime
}

func decrementExpiresIn() {
	for {
		time.Sleep(time.Second)
		expiresIn--
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	refreshToken(e)
	listSwitches(e, ch)
	listAccessPoints(e, ch)
	listMobilityControllers(e, ch)
	listTopClients(e, ch)
	listSites(e, ch)

}

func init() {

	flag.BoolVar(&verbose, "v", false, "Enable verbose mode")
	flag.StringVar(&configFile, "f", "exporter_config.yaml", "Specify config file")

	flag.Usage = func() {
		fmt.Println("Usage: aruba_exporter [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
}

func main() {

	fmt.Println(time.Now().Format(time.RFC3339), "Aruba Central Exporter v", version, " is running...")

	flag.Parse()

	go decrementExpiresIn()

	config := Config{}
	response := Response{}
	authenticate(&config, configFile, &response)

	arubaEndpoint := config.ArubaEndpoint
	arubaAccessToken := response.AccessToken
	arubaRefreshToken := response.RefreshToken
	exporterEndpoint := config.ExporterConfig[0].ExporterEndpoint
	exporterPort := config.ExporterConfig[1].ExporterPort

	exporter := NewExporter(arubaEndpoint, arubaAccessToken, arubaRefreshToken)
	prometheus.MustRegister(exporter)

	http.Handle(exporterEndpoint, promhttp.Handler())

	fmt.Println(time.Now().Format(time.RFC3339), "Server listening on port", exporterPort)
	err := http.ListenAndServe(exporterPort, nil)

	if err != nil {
		if err.Error() == "listen tcp :8080: bind: address already in use" {
			fmt.Println("Error: Port", exporterPort, "is already in use.")
		} else {
			fmt.Println("Error starting server:", err)
		}
	}

}

func listAccessPoints(e *Exporter, ch chan<- prometheus.Metric) {

	url := e.arubaEndpoint + "monitoring/v2/aps?calculate_total=true&calculate_client_count=true&calculate_ssid_count=true&show_resource_details=true"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+e.arubaAccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	// Parse JSON
	var apResponse ApResponse
	if err := json.Unmarshal(body, &apResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, a := range apResponse.AccessPoints {

		ch <- prometheus.MustNewConstMetric(apClientCount, prometheus.GaugeValue, float64(a.ClientCount), a.Name, a.GroupName, a.Site, a.Status, a.FirmwareVersion, a.Model)
		ch <- prometheus.MustNewConstMetric(apCpuUtilization, prometheus.GaugeValue, float64(a.CpuUtilization), a.Name, a.GroupName, a.Site, a.Status, a.FirmwareVersion, a.Model)
		ch <- prometheus.MustNewConstMetric(apMemFree, prometheus.GaugeValue, float64(a.MemFree), a.Name, a.GroupName, a.Site, a.Status, a.FirmwareVersion, a.Model)
		ch <- prometheus.MustNewConstMetric(apMemTotal, prometheus.GaugeValue, float64(a.MemTotal), a.Name, a.GroupName, a.Site, a.Status, a.FirmwareVersion, a.Model)
		ch <- prometheus.MustNewConstMetric(apUptime, prometheus.GaugeValue, float64(a.Uptime), a.Name, a.GroupName, a.Site, a.Status, a.FirmwareVersion, a.Model)

		for _, r := range a.Radios {

			ch <- prometheus.MustNewConstMetric(apRadioTxPower, prometheus.GaugeValue, float64(r.TxPower), strconv.Itoa(r.Band), r.Channel, r.RadioName, a.Name)
			ch <- prometheus.MustNewConstMetric(apRadioUtilization, prometheus.GaugeValue, float64(r.Utilization), strconv.Itoa(r.Band), r.Channel, r.RadioName, a.Name)
		}
	}
	if verbose {

		fmt.Println("\nmonitoring/v2/aps - HTTP Status Code:", resp.StatusCode)

		for key, value := range resp.Header {
			fmt.Printf(" (%s: %s),", key, value)
		}
	}

}

func listMobilityControllers(e *Exporter, ch chan<- prometheus.Metric) {

	url := e.arubaEndpoint + "monitoring/v1/mobility_controllers?calculate_total=false"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+e.arubaAccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	// Parse JSON
	var mcResponse McResponse
	if err := json.Unmarshal(body, &mcResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, m := range mcResponse.MobilityControllers {

		ch <- prometheus.MustNewConstMetric(mcCpuUtilization, prometheus.GaugeValue, float64(m.CpuUtilization), m.Name, m.GroupName, m.Mode, m.Model, m.Site, m.Status, m.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(mcMemFree, prometheus.GaugeValue, float64(m.MemFree), m.Name, m.GroupName, m.Mode, m.Model, m.Site, m.Status, m.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(mcMemTotal, prometheus.GaugeValue, float64(m.MemTotal), m.Name, m.GroupName, m.Mode, m.Model, m.Site, m.Status, m.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(mcUptime, prometheus.GaugeValue, float64(m.Uptime), m.Name, m.GroupName, m.Mode, m.Model, m.Site, m.Status, m.FirmwareVersion)
	}

	if verbose {
		fmt.Println("\nmonitoring/v1/mobility_controllers - HTTP Status Code:", resp.StatusCode)
		for key, value := range resp.Header {
			fmt.Printf(" (%s: %s),", key, value)
		}
	}

}

func listSites(e *Exporter, ch chan<- prometheus.Metric) {

	url := e.arubaEndpoint + "branchhealth/v1/site?limit=100&column=device_total&order=desc"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+e.arubaAccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	// Parse JSON

	var siteResponse SiteResponse

	if err := json.Unmarshal(body, &siteResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, s := range siteResponse.Sites {

		ch <- prometheus.MustNewConstMetric(siteConnectedCount, prometheus.GaugeValue, float64(s.ConnectedCount), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteDeviceDown, prometheus.GaugeValue, float64(s.DeviceDown), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteDeviceHighCh24, prometheus.GaugeValue, float64(s.DeviceHighCh24), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteDeviceHighCh5, prometheus.GaugeValue, float64(s.DeviceHighCh5), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteDeviceHighCpu, prometheus.GaugeValue, float64(s.DeviceHighCpu), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteDeviceHighMem, prometheus.GaugeValue, float64(s.DeviceHighMem), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteDeviceHighNoise24, prometheus.GaugeValue, float64(s.DeviceHighNoise24), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteDeviceHighNoise5, prometheus.GaugeValue, float64(s.DeviceHighNoise5), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteDeviceUp, prometheus.GaugeValue, float64(s.DeviceUp), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteWiredCpuHigh, prometheus.GaugeValue, float64(s.WiredCPUHigh), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteWiredDeviceStatusDown, prometheus.GaugeValue, float64(s.WiredDeviceSatusDown), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteWiredDeviceStatusUp, prometheus.GaugeValue, float64(s.WiredDeviceStatusUp), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteWiredMemHigh, prometheus.GaugeValue, float64(s.WiredMemHigh), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteWlanCpuHigh, prometheus.GaugeValue, float64(s.WiredCPUHigh), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteWlanDeviceStatusDown, prometheus.GaugeValue, float64(s.WlanDeviceStatusDown), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteWlanDeviceStatusUp, prometheus.GaugeValue, float64(s.WlanDeviceStatusDown), s.Name, s.Id)
		ch <- prometheus.MustNewConstMetric(siteWlanMemHigh, prometheus.GaugeValue, float64(s.WlanMemHigh), s.Name, s.Id)

	}

	if verbose {
		fmt.Println("\nbranchhealth/v1/site - HTTP Status Code:", resp.StatusCode)
		for key, value := range resp.Header {
			fmt.Printf(" (%s: %s),", key, value)
		}
	}

}

func listSwitches(e *Exporter, ch chan<- prometheus.Metric) {

	url := e.arubaEndpoint + "monitoring/v1/switches?show_resource_details=true&calculate_client_count=true"

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+e.arubaAccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	// Parse JSON
	var switchResponse SwitchResponse
	if err := json.Unmarshal(body, &switchResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, s := range switchResponse.Switches {

		ch <- prometheus.MustNewConstMetric(switchClientCount, prometheus.GaugeValue, float64(s.ClientCount), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion, s.Model)
		ch <- prometheus.MustNewConstMetric(switchCpuUtilization, prometheus.GaugeValue, float64(s.CPUUtilization), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion, s.Model)
		ch <- prometheus.MustNewConstMetric(switchMemFree, prometheus.GaugeValue, float64(s.ClientCount), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion, s.Model)
		ch <- prometheus.MustNewConstMetric(switchMemTotal, prometheus.GaugeValue, float64(s.ClientCount), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion, s.Model)
		ch <- prometheus.MustNewConstMetric(switchUsage, prometheus.GaugeValue, float64(s.Usage), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion, s.Model)
		ch <- prometheus.MustNewConstMetric(switchUptime, prometheus.GaugeValue, float64(s.Uptime), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion, s.Model)
	}

	if verbose {
		fmt.Println("\nmonitoring/v1/switches - HTTP Status Code:", resp.StatusCode)

		for key, value := range resp.Header {
			fmt.Printf(" (%s: %s),", key, value)
		}
	}

}

func listTopClients(e *Exporter, ch chan<- prometheus.Metric) {

	url := e.arubaEndpoint + "monitoring/v1/clients/bandwidth_usage/topn?count=100"

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+e.arubaAccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	// Parse JSON
	var topNClientResponse TopNClientResponse
	if err := json.Unmarshal(body, &topNClientResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, t := range topNClientResponse.Clients {

		ch <- prometheus.MustNewConstMetric(clientRxDataBytes, prometheus.GaugeValue, float64(t.RxDataBytes), t.Name, t.MacAddress)
		ch <- prometheus.MustNewConstMetric(clientTxDataBytes, prometheus.GaugeValue, float64(t.TxDataBytes), t.Name, t.MacAddress)

	}

	if verbose {
		fmt.Println("\nmonitoring/v1/clients - HTTP Status Code:", resp.StatusCode)

		for key, value := range resp.Header {
			fmt.Printf(" (%s: %s),", key, value)
		}
	}
}

func refreshToken(e *Exporter) {

	if expiresIn < 60 {

		config := Config{}
		readConfig(&config, configFile)

		clientId := config.ArubaApplicationCredentials[0].ClientID
		clientSecret := config.ArubaApplicationCredentials[1].ClientSecret

		url := e.arubaEndpoint + "oauth2/token?client_id=" + clientId + "&client_secret=" + clientSecret + "&grant_type=refresh_token&refresh_token=" + e.arubaRefreshToken

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)

		}

		req.Header.Set("Authorization", "Bearer "+e.arubaAccessToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)

		}

		// Parse JSON

		var tokenResponse TokenResponse

		if err := json.Unmarshal(body, &tokenResponse); err != nil {
			fmt.Println("Error parsing JSON:", err)

		}

		expiresIn = tokenResponse.ExpiresIn

		e.arubaAccessToken = tokenResponse.AccessToken
		e.arubaRefreshToken = tokenResponse.RefreshToken
	}

}
