from flask import Flask
import os, logging

from src.routes.inventory import inventory_bp
from src.routes.reports import reports_bp


def create_app():
    app = Flask(__name__)

    app.config["SECRET_KEY"] = "super-secret-key-hardcoded"  # noqa: S105 — Intentional: hardcoded secret
    app.config["DEBUG"] = True  # Intentional: debug mode enabled

    app.register_blueprint(inventory_bp, url_prefix="/api/inventory")
    app.register_blueprint(reports_bp, url_prefix="/api/reports")

    @app.route("/")
    def index():
        return {"status": "ok", "app": "cq-demo-app-002", "version": "1.0.0"}

    @app.route("/health")
    def health():
        return {"status": "healthy"}

    return app


app = create_app()

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000, debug=True)
