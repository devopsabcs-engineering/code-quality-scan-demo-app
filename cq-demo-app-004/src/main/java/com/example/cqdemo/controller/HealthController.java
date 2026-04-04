package com.example.cqdemo.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@RestController
public class HealthController {

    @GetMapping("/")
    public Map<String, String> root() {
        return Map.of(
            "status", "ok",
            "app", "cq-demo-app-004",
            "framework", "Spring Boot"
        );
    }
}
