pipeline {
    agent any

    stages {
        stage('Update Components') {
            when {
                anyOf {
                    branch 'master'
                    }
                }
            steps {
                echo "Updating"
                sh "docker pull golang:1.15-alpine"
            }
        }
        stage('Build') {
            when {
                anyOf {
                    branch 'master'
                    }
                }
            steps {
                echo "Building"
                sh "docker build -t localhost:5000/ystv/web-api:$BUILD_ID ."
            }
        }
        stage('Upload & Cleanup') {
            when {
                anyOf {
                    branch 'master'
                    }
                }
            steps {
                echo "Uploading To Registry"
                sh "docker push localhost:5000/ystv/web-api:$BUILD_ID" // Uploaded to registry
                echo "Performing Cleanup"
                sh "docker image prune -f --filter label=site=api --filter label=stage=builder" // Removing the local builder image
                sh "docker image rm localhost:5000/ystv/web-api:$BUILD_ID" // Removing the local builder image
            }
        }
        stage('Deploy') {
            when {
                anyOf {
                    branch 'master'
                    }
                }
            steps {
                echo "Deploying"
                sh "docker pull localhost:5000/ystv/web-api:$BUILD_ID" // Pulling image from local registry
                script {
                    try {
                        sh "docker kill ystv-dev-api" // Stop old container
                    }
                    catch (err) {
                        echo "Couldn't find container to stop"
                        echo err.getMessage()
                    }
                }
                sh "docker run -d --rm -p 1336:8081 --env-file /YSTV-ENVVARS/api.env --name ystv-dev-api localhost:5000/ystv/web-api:$BUILD_ID" // Deploying site
                sh 'docker image prune -a -f --filter "label=site=api"' // remove old image
            }
        }
    }
    post {
        success {
            echo 'Very cash-money'
        }
        failure {
            echo 'That is not ideal, cheeky bugger'
        }
    }
}