resource "aws_security_group" "ingress-all-pennyearned-database" {
  name   = "allow-all-sg"
  vpc_id = aws_vpc.pennyearned-database-vpc.id
  ingress {
    cidr_blocks = [
      "0.0.0.0/0"
    ]
    from_port = 22
    to_port   = 22
    protocol  = "tcp"
  }
  // Terraform removes the default rule
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
resource "aws_security_group" "ingress-postgres-database" {
  name        = "allow-all-postgres"
  vpc_id      = aws_vpc.pennyearned-database-vpc.id
  description = "The rule to allow connections to the db"
  ingress {
    cidr_blocks = [
      "0.0.0.0/0"
    ]
    from_port = 5432
    to_port   = 5432
    protocol  = "tcp"
  }
  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
}

resource "aws_security_group" "ingress-redis-database" {
  name        = "allow-all-redis"
  vpc_id      = aws_vpc.pennyearned-database-vpc.id
  description = "The rule to allow connections to redis"
  ingress {
    cidr_blocks = [
      "0.0.0.0/0"
    ]
    from_port = 6379
    to_port   = 6379
    protocol  = "tcp"
  }
  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
}
