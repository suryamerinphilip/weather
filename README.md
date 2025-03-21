# GCP-Cloud-Run-Golang-App

In this project, I've deployed a containerized Weather Golang application on Google Cloud Run, using Terraform for infrastructure management and Continuous Deployment through integration with GitHub and Snyk for security. The deployment uses dual Cloud Build pipelines; one automates Terraform infrastructure deployment, and the other manages Docker image updates. 

The application leverages the Gin HTTP framework for efficient HTML delivery and uses CSS, JavaScript, and APIs like OpenCage and OpenWeather for global weather information retrieval. Firebase CDN Hosting has been added for caching, ensuring an impressively swift response time for users globally. To enhance user experience, a custom domain name from Google Domains has been assigned, providing a user-friendly web address. This robust strategy ensures a scalable, secure, and efficient application deployment on the Cloud Run platform.

Link to Website: https://weather-app.reggietestgcpdomain.com/

## Architecture Breakdown

The application is broken down into the architecture below:

![applications](https://github.com/rjones18/Images/blob/main/Cloud%20Run%20Application.png)
