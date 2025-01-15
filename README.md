# Image Uploader POC

This project is a Proof of Concept (POC) for building web applications using the following stack:

- **Golang**: Backend server for handling HTTP requests and managing file uploads.
- **Templ**: Templating engine for rendering dynamic HTML on the server side.
- **AlpineJS**: Lightweight frontend framework for adding interactivity without heavy dependencies.
- **Vite**: Fast and modern build tool for bundling and optimizing assets.

## Features

- Drag-and-drop image upload support.
- File browsing for uploading images from computers or mobile devices.
- Responsive design to ensure compatibility with various devices.
- Server-side validation for image uploads.
- Clear feedback for successful and failed uploads.

## Objectives

The primary goal of this project is to demonstrate the viability of using Golang, Templ, AlpineJS, and Vite to build a modern web application. The stack highlights simplicity, performance, and security.

## Getting Started

### Prerequisites

Make sure you have the following installed:

- [Golang](https://golang.org/) (version 1.19 or later recommended)
- [Node.js](https://nodejs.org/) (for Vite, version 16 or later recommended)
- [npm](https://www.npmjs.com/) or [yarn](https://yarnpkg.com/)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/rubensdev/go-image-uploader.git
cd go-image-uploader
```

2. Install frontend dependencies:

```bash
cd frontend
npm install
```

3. Build the frontend assets:

```bash
npm run build
```

4. Build the executable:

```bash
cd ..
make build
```

5. Run the executable in production mode:
  
```bash
cd bin
./app -env=production
```

6. Open your browser and navigate to `http://localhost:8080` to view the application.


## Usage

1. Drag and drop image files into the upload area, or click to browse files from your device.
2. The application will validate and upload the images to the server.
3. Successfully uploaded images are stored in the `uploads/` directory on the server.

## Future Improvements

I'm sure there is a lot of room for improvement, but I focused on getting a working setup for this stack =).

## License

This project is licensed under the MIT License. See the [https://mit-license.org/](LICENSE) file for details.

## Acknowledgments

Special thanks to the developers of Golang, Templ, AlpineJS, and Vite for their excellent tools that made this project possible.

