package main

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
)

var appsFile = appsDir + "/apps.json"
var appsDir = "./test"

func setup() {
	err := os.Mkdir(appsDir, 0755)
	if err != nil {
		log.Fatalf("Error on setup %s", err.Error())
	}
}

func teardown() {
	err := os.RemoveAll(appsDir)
	if err != nil {
		log.Fatalf("Error on teardown %s", err.Error())
	}
}

func Test_parseLabels(t *testing.T) {

	type args struct {
		containers []types.Container
	}

	tests := []struct {
		name    string
		args    args
		want    Domains
		wantErr bool
	}{
		{
			name: "Assert labels are parsed correctly into App struct",
			args: args{
				containers: []types.Container{
					{
						Names: []string{"MyApp", "MyApp"},
						Labels: map[string]string{
							"pihole.domain.name": "MyApp.test.xyz",
							"pihole.domain.ip":   "1.2.3.4",
						},
					},
				},
			},
			want: Domains{
				{
					Name: "MyApp.test.xyz",
					IP:   "1.2.3.4",
				},
			},
			wantErr: false,
		},
		{
			name: "Assert unwanted labels are ignored",
			args: args{
				containers: []types.Container{
					{
						Names: []string{"NoLabelApp", "NoLabelApp"},
						Labels: map[string]string{
							"some.random.label": "some.random.value",
						},
					},
				},
			},
			want:    Domains{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPiholeLabels(tt.args.containers)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLabels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLabels() = %v, want %v", got, tt.want)
			}
		})
	}
}
