#version 460 core

layout (location = 0) in vec3 aPos;

out vec4 vertexColor; // specify color output for frag shader

void main() {
    gl_Position = vec4(aPos, 1.0);
    vertexColor = vec4(aPos, 1);
}
