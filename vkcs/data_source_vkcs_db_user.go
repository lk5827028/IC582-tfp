package vkcs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatabaseUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"dbms_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"host": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"databases": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating vkcs database client: %s", err)
	}

	id := d.Get("id").(string)
	userID := strings.SplitN(id, "/", 2)
	if len(userID) != 2 {
		return diag.Errorf("invalid vkcs_db_user id: %s", id)
	}

	dbmsID := userID[0]
	userName := userID[1]
	dbmsResp, err := getDBMSResource(DatabaseV1Client, dbmsID)
	if err != nil {
		return diag.Errorf("error while getting resource: %s", err)
	}
	var dbmsType string
	if _, ok := dbmsResp.(instanceResp); ok {
		dbmsType = dbmsTypeInstance
	}
	if _, ok := dbmsResp.(dbClusterResp); ok {
		dbmsType = dbmsTypeCluster
	}
	exists, userObj, err := databaseUserExists(DatabaseV1Client, dbmsID, userName, dbmsType)
	if err != nil {
		return diag.Errorf("error checking if vkcs_db_user %s exists: %s", d.Id(), err)
	}

	if !exists {
		d.SetId("")
		return nil
	}

	d.SetId(id)
	d.Set("name", userName)

	databases := flattenDatabaseUserDatabases(userObj.Databases)
	if err := d.Set("databases", databases); err != nil {
		return diag.Errorf("unable to set databases: %s", err)
	}
	d.Set("dbms_id", dbmsID)
	return nil
}