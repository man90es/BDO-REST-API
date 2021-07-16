package models

type MaintenanceError struct{}

func (e *MaintenanceError) Error() string {
	return "BDO servers are currently under maintenance"
}

func (e *MaintenanceError) HTTPCode() int {
	return 503
}
