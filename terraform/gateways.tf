resource "aws_internet_gateway" "pennyearned-database-gw" {
  vpc_id = aws_vpc.pennyearned-database-vpc.id
  tags = {
    Name = "pennyearned-database-gw"
  }
}
