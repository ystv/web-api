String registryEndpoint = 'registry.comp.ystv.co.uk'

def image
String imageName = "ystv/web-api:${env.BRANCH_NAME}-${env.BUILD_ID}"

pipeline {
  agent {
    label 'docker'
  }

  environment {
    DOCKER_BUILDKIT = '1'
  }

  stages {
    stage('Build image') {
      steps {
        script {
          docker.withRegistry('https://' + registryEndpoint, 'docker-registry') {
            image = docker.build(imageName)
          }
        }
      }
    }

    stage('Push image to registry') {
      steps {
        script {
          docker.withRegistry('https://' + registryEndpoint, 'docker-registry') {
            image.push()
            if (env.BRANCH_IS_PRIMARY) {
              image.push('latest')
            }
          }
        }
      }
    }

    stage('Deploy') {
      when {
        expression { env.BRANCH_IS_PRIMARY }
      }
      steps {
        build(job: 'Deploy Nomad Job', parameters: [
         string(name: 'JOB_FILE', value: 'web-api-dev.nomad'),
         text(name: 'TAG_REPLACEMENTS', value: "${registryEndpoint}/${imageName}")
        ])
      }
    }
  }
}
