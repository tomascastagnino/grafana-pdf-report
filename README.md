
# Grafana PDF Reporter

Grafana PDF Reporter is a Go-based service that generates PDF reports from Grafana dashboards. It provides an endpoint to fetch dashboards, panels, and their associated images, rendering them into a single PDF document.

## Features

- Fetch and render Grafana dashboards into PDF.
- Automatically handle Grafana panel images.
- Support for exporting and downloading dashboards as PDF files.
- Customizable dashboard and panel layouts.

## Requirements

- **Grafana**: Ensure you have Grafana running and accessible.
- **Grafana Image Renderer Plugin**: This plugin is necessary to render panels as images. You can find the plugin [here](https://grafana.com/grafana/plugins/grafana-image-renderer).

## Setup

The application is designed to be run using Docker. You should change the `docker-compose.yml` to point to your own Grafana instance. However, you can use the provided `docker-compose.yml` for testing purposes. Below are the steps to set up the application using Docker.

### Docker

1. **Clone the Repository**

   ```bash
   git clone https://github.com/your-username/grafana-pdf-reporter.git
   cd grafana-pdf-reporter
   ```

2. **Build and Run the Docker Containers**

   Ensure Docker is installed and running on your system. The `docker-compose.yml` file provided in the repository is configured to run the application along with Grafana and the Grafana Image Renderer plugin.

   To build and run the Docker containers:

   ```bash
   docker-compose up --build
   ```

   This command will start the following services:

   - **Grafana**: Accessible at `http://localhost:3000`
   - **Grafana Image Renderer**: Runs in the background to render images of Grafana panels.
   - **Grafana PDF Reporter**: The service itself, accessible at `http://localhost:9090`.

   You can use this setup for testing and development purposes. Adjustments to the `docker-compose.yml` file might be necessary depending on your production environment.

## Usage

Once the application is running, you can access the Grafana PDF Reporter service via `http://localhost:9090`. The main functionalities include:

- **Fetching Dashboards**: The service will interact with Grafana to retrieve dashboards and render them as PDFs.
- **Home Page**: A simple interface to list all available dashboards.
- **PDF Export**: Export selected dashboards as PDF documents.

## Endpoints

- **Home Page**: `GET /`
- **Dashboard List**: `GET /api/v1/dashboards`
- **Fetch Dashboard**: `GET /api/v1/dashboard/{dashboard_id}`
- **Fetch Panel Image**: `GET /api/v1/dashboard/{dashboard_id}/panel/{panel_id}`
- **Export as PDF**: Accessible via the web interface.

## Contributing

Feel free to submit issues, fork the repository, and send pull requests. Contributions are welcome!

## License

This project is licensed under the MIT License.
