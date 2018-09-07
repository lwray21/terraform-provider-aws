package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAwsServerlessRepositoryApplication_Basic(t *testing.T) {
	var stack cloudformation.Stack
	stackName := fmt.Sprintf("tf-acc-test-basic-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCloudFormationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsServerlessRepositoryApplicationConfig(stackName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckerverlessRepositoryApplicationExists("aws_serverlessrepository_application.postgres-rotator", &stack),
				),
			},
		},
	})
}

func testAccAwsServerlessRepositoryApplicationConfig(stackName string) string {
	return fmt.Sprintf(`
resource "aws_serverlessrepository_application" "postgres-rotator" {
  name           = "%[1]s"
  application_id = "arn:aws:serverlessrepo:us-east-1:297356227824:applications/SecretsManagerRDSPostgreSQLRotationSingleUser"
  parameters = {
    functionName = "func-%[1]s"
    endpoint     = "secretsmanager.us-east-2.amazonaws.com"
  }
}`, stackName)
}

func testAccCheckerverlessRepositoryApplicationExists(n string, stack *cloudformation.Stack) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*AWSClient).cfconn
		params := &cloudformation.DescribeStacksInput{
			StackName: aws.String(rs.Primary.ID),
		}
		resp, err := conn.DescribeStacks(params)
		if err != nil {
			return err
		}
		if len(resp.Stacks) == 0 {
			return fmt.Errorf("CloudFormation stack not found")
		}

		return nil
	}
}