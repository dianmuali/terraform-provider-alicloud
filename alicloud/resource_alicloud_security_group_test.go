package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_security_group", &resource.Sweeper{
		Name: "alicloud_security_group",
		F:    testSweepSecurityGroups,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_instance",
		},
	})
}

func testSweepSecurityGroups(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var groups []ecs.SecurityGroup
	req := ecs.CreateDescribeSecurityGroupsRequest()
	req.RegionId = conn.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		resp, err := conn.ecsconn.DescribeSecurityGroups(req)
		if err != nil {
			return fmt.Errorf("Error retrieving Security Groups: %s", err)
		}
		if resp == nil || len(resp.SecurityGroups.SecurityGroup) < 1 {
			break
		}
		groups = append(groups, resp.SecurityGroups.SecurityGroup...)

		if len(resp.SecurityGroups.SecurityGroup) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, v := range groups {
		name := v.SecurityGroupName
		id := v.SecurityGroupId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Security Group: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Security Group: %s (%s)", name, id)
		req := ecs.CreateDeleteSecurityGroupRequest()
		req.SecurityGroupId = id
		if _, err := conn.ecsconn.DeleteSecurityGroup(req); err != nil {
			log.Printf("[ERROR] Failed to delete Security Group (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlicloudSecurityGroup_basic(t *testing.T) {
	var sg ecs.DescribeSecurityGroupAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(
						"alicloud_security_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"alicloud_security_group.foo",
						"name",
						"tf-testAccSecurityGroupConfig"),
				),
			},
		},
	})

}

func TestAccAlicloudSecurityGroup_withVpc(t *testing.T) {
	var sg ecs.DescribeSecurityGroupAttributeResponse
	var vpc vpc.DescribeVpcAttributeResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_security_group.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupConfig_withVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(
						"alicloud_security_group.foo", &sg),
					testAccCheckVpcExists(
						"alicloud_vpc.vpc", &vpc),
					resource.TestCheckResourceAttr(
						"alicloud_security_group.foo",
						"inner_access",
						"true"),
				),
			},
		},
	})

}

func testAccCheckSecurityGroupExists(n string, sg *ecs.DescribeSecurityGroupAttributeResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SecurityGroup ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		d, err := client.DescribeSecurityGroupAttribute(rs.Primary.ID)

		log.Printf("[WARN] security group id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}
		if d.SecurityGroupId == rs.Primary.ID {
			*sg = d
		}
		return nil
	}
}

func testAccCheckSecurityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_security_group" {
			continue
		}

		group, err := client.DescribeSecurityGroupAttribute(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) || IsExceptedErrors(err, []string{InvalidSecurityGroupIdNotFound}) {
				continue
			}
			return err
		}

		if group.SecurityGroupId != "" {
			return fmt.Errorf("Error SecurityGroup still exist")
		}
	}
	return nil
}

const testAccSecurityGroupConfig = `
resource "alicloud_security_group" "foo" {
  name = "tf-testAccSecurityGroupConfig"
}
`

const testAccSecurityGroupConfig_withVpc = `
variable "name" {
  default = "tf-testAccSecurityGroupConfig_withVpc"
}
resource "alicloud_security_group" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_vpc" "vpc" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}
`
