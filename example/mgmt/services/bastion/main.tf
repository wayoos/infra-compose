// Simple terraform file as infra-compose example.

resource "local_file" "demo" {
    content     = "Bastion applied"
    filename = "${path.module}/applied-demo.txt"
}