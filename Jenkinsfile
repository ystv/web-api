pipeline {
    agent any

    environment {
        REGISTRY_ENDPOINT = credentials('docker-registry-endpoint')
    }

    stages {
        stage('Update Components') {
            steps {
                sh "docker pull golang:1.16-alpine" // Update with current Go image
            }
        }
        stage('Build') {
            steps {
                sh 'docker build -t $REGISTRY_ENDPOINT/ystv/web-api:$BUILD_ID .'
            }
        }
        stage('Registry Upload') {
            steps {
                sh 'docker push $REGISTRY_ENDPOINT/ystv/web-api:$BUILD_ID' // Uploaded to registry
            }
        }
        stage('Deploy') {
            stages {
                stage('Staging') {
                    when {
                        branch 'master'
                        not {
                            tag pattern: "^v(?P<major>0|[1-9]\\d*)\\.(?P<minor>0|[1-9]\\d*)\\.(?P<patch>0|[1-9]\\d*)", comparator: "REGEXP"
                        }
                    }
                    environment {
                        APP_ENV = credentials('wapi-staging-env')
                        TARGET_SERVER = credentials('staging-server-address')
                        TARGET_PATH = credentials('staging-server-path')
                    }
                    steps {
                        sshagent(credentials : ['staging-server-key']) {
                            script {
                                sh 'rsync -av $APP_ENV deploy@$TARGET_SERVER:$TARGET_PATH/web-api/.env'
                                sh '''ssh -tt deploy@$TARGET_SERVER << EOF
                                    docker pull $REGISTRY_ENDPOINT/ystv/web-api:$BUILD_ID
                                    docker rm -f ystv-web-api
                                    docker run -d -p 1336:8081 --env-file $TARGET_PATH/web-api/.env --name ystv-web-api $REGISTRY_ENDPOINT/ystv/web-api:$BUILD_ID
                                    docker image prune -a -f --filter "label=site=api"
                                    exit 0
                                EOF'''
                            }
                        }
                    }
                }
                stage('Production') {
                    when {
                        tag pattern: "^v(?P<major>0|[1-9]\\d*)\\.(?P<minor>0|[1-9]\\d*)\\.(?P<patch>0|[1-9]\\d*)", comparator: "REGEXP" // Checking if it is main semantic version release
                    }
                    environment {
                        APP_ENV = credentials('wapi-prod-env')
                        TARGET_SERVER = credentials('prod-server-address')
                        TARGET_PATH = credentials('prod-server-path')
                    }
                    steps {
                        sshagent(credentials : ['prod-server-key']) {
                            script {
                                sh 'rsync -av $APP_ENV deploy@$TARGET_SERVER:$TARGET_PATH/web-api/.env'
                                sh '''ssh -tt deploy@$TARGET_SERVER << EOF
                                    docker pull $REGISTRY_ENDPOINT/ystv/web-api:$BUILD_ID
                                    docker rm -f ystv-web-api
                                    docker run -d -p 1336:8081 --env-file $TARGET_PATH/web-api/.env --name ystv-web-api $REGISTRY_ENDPOINT/ystv/web-api:$BUILD_ID
                                    docker image prune -a -f --filter "label=site=api"
                                    exit 0
                                EOF'''
                            }
                        }
                    }
                }
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
        always {
            sh "docker image prune -f --filter label=site=api --filter label=stage=builder" // Removing the local builder image
            sh 'docker image prune -a -f --filter "label=site=api"' // remove old image
        }
    }
}