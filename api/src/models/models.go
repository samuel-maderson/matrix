package models

type InstancePayload struct {
	InstanceName         string `json:"instance_name"`
	SSHKeyName           string `json:"ssh_key_name"`
	NumberOfInstances    int    `json:"number_of_instances"`
	InstanceType         string `json:"instance_type"`
	Launched             string `json:"launched"`
	WhenInstanceStart    string `json:"when_instance_start"`
	WhenInstanceShutdown string `json:"when_instance_shutdown"`
}

type RequestData struct {
	InstanceName         string `json:"instance_name"`
	SSHKeyName           string `json:"ssh_key_name"`
	NumberOfInstances    string `json:"number_of_instances"`
	InstanceType         string `json:"instance_type"`
	Launched             string `json:"launched"`
	WhenInstanceStart    string `json:"when_instance_start"`
	WhenInstanceShutdown string `json:"when_instance_shutdown"`
}
