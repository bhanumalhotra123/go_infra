package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		sgArgs := &ec2.SecurityGroupArgs{
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(8080),
					ToPort:     pulumi.Int(8080),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		}

		sg, err := ec2.NewSecurityGroup(ctx, "jenkins-sg", sgArgs)
		if err != nil {
			return err
		}

		kp, err := ec2.NewKeyPair(ctx, "local-ssh", &ec2.KeyPairArgs{PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDGSeYq+n1KAbB2DPFbY58UmvZo89awdGwKfJCytVrsV7jnCoMICQ4LvRuTyM7ZeJPK/os+cv3aYTQesVCZpyotipAQmQE57iiXDVykgGtFkzWnHDp+pGp8aWtqHhHbfgTWg+gXe7nwx7pqtlO+mDdQ28Rd6vyOMbt/x+xD23LOS4wyxPsBSKfDOmgyq7515a4AZRfqOp0yPR1OwB7QiEHoh+0PrRVvCrEbmjYttuOMLGy5TPj191eIzvt2FjAP3rmZn1NN2eqFsDUhsrpgm5MJDhAHaTrW0ZJ6eKw3P8K5flneB8iuFAzBbclcHRYKoFndvZpaMFOn3eq9IzTNKcfCkwKX0V3uv8NcQ2PA4rGPuRzzMD6bSgTnzNRwAGGG4ka4NMv7aC9WDxNZc7k9Rjw9whpjJjzKEFogxTYPai3wnkPk2sHy3QtFUlLFc7RouS77afEVEwM5lHvC+gcXoVJKOXGEQOtWD+KVOjlmImKrdLsFzA6IP2j0vF7q/sNNyOs= HP@bhanumalhotra")})
		if err != nil {
			return err
		}

		jenkinsServer, err := ec2.NewInstance(ctx, "jenkins-server", &ec2.InstanceArgs{
			InstanceType:        pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			Ami:                 pulumi.String("ami-04cb4ca688797756f"),
			KeyName:             kp.KeyName,
		})

		if err != nil {
			return err
		}

		fmt.Println(jenkinsServer.PublicIp)
		fmt.Println(jenkinsServer.PublicDns)

		ctx.Export("publicIp", jenkinsServer.PublicIp)
		ctx.Export("publicHostName", jenkinsServer.PublicDns)
		return nil
	})
}
