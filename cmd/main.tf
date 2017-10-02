// Simple terraform file as infra-compose example.

resource "local_file" "demo" {
    content     = "Global applied"
    filename = "${path.module}/applied-demo.txt"
}