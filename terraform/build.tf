terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}
provider "aws" {
  region = "us-east-2"
}
//resource "aws_key_pair" "pennyearned_database_key" {
//  key_name   = "pennyearned_database_key"
//  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCDvv55uxvErvduaC4VFTmVzornxTgVDen+x3ot2tpawdjgieGgJG4tc0Gl729CColfERHHuSw5GrYmLe9chZArufqMx/6t0PCNDAGUrK3F43ironHN1f5Ht9T3U3ABWJCBIgU2we7GK4OS46HyB/wrw1k1AH67DlZgA5AL/Xs7cTzPirjZxbPPaOAC9CpSivaIBffjiLWYkjeu3VIEQ0F0i13u/40d381qzk/XTRj6FJGLDQIXgX4xbJexhBzsEHFbH2Giz4og3tlH9shrzDogxJJrCZP1WHk7GuveE/mT99+tHV8oUkQ+jaUbebaez7N8oM3AIgV0DR6OqycnFTuX"
//}
resource "aws_instance" "pennyearned_database" {
  ami           = "ami-02238ac43d6385ab3"
  instance_type = "t3.nano"
  //key_name      = aws_key_pair.pennyearned_database_key.key_name
  key_name = "pennyearned_database_key"
  tags = {
    Name = "PennyearnedDatabaseInstace"
  }
  security_groups = [
    "${aws_security_group.ingress-all-pennyearned-database.id}",
    "${aws_security_group.ingress-postgres-database.id}",
    "${aws_security_group.ingress-redis-database.id}"
  ]
  subnet_id = aws_subnet.pennyearned-database-subnet.id
  root_block_device {
    volume_size           = "15"
    volume_type           = "gp2"
    delete_on_termination = false
  }
}
