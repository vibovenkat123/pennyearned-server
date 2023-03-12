resource "aws_subnet" "pennyearned-database-subnet" {
  cidr_block              = cidrsubnet(aws_vpc.pennyearned-database-vpc.cidr_block, 3, 1)
  vpc_id                  = aws_vpc.pennyearned-database-vpc.id
  availability_zone       = "us-east-2a"
  map_public_ip_on_launch = true
}
resource "aws_route_table" "route-table-pennyearned-database" {
  vpc_id = aws_vpc.pennyearned-database-vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.pennyearned-database-gw.id
  }
  tags = {
    Name = "route-table-pennyearned-database"
  }
}
resource "aws_route_table_association" "pennyearned-database-subnet-association" {
  subnet_id      = aws_subnet.pennyearned-database-subnet.id
  route_table_id = aws_route_table.route-table-pennyearned-database.id
}
