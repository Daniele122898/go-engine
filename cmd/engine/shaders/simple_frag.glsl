#version 460

out vec4 FragColor;

in vec4 vertexColor; // Input from vert shader. Same name and type

void main() {
    FragColor = vertexColor;
}
