pipeline {
  agent any
  environment {
    DOCKER_REGISTRY = 'registry.cn-hangzhou.aliyuncs.com'
  }
  stages {
    stage("检出代码") {
      steps {
        checkout([
          $class: 'GitSCM',
          branches: [[name: env.GIT_BUILD_REF]],
          userRemoteConfigs: [[
            url: env.GIT_REPO_URL,
            credentialsId: env.CREDENTIALS_ID
        ]]])
      }
    }
    stage('构建镜像') {
      when {
        buildingTag()
      }
      steps {
        script {
          try {
            IMAGE_VERSION = "${env.TAG_NAME}"
            IMAGE_NAME = "${DEPOT_NAME}"
            IMAGE_FULL_NAME = "${DOCKER_REGISTRY}/dongfg/${IMAGE_NAME}"
            withCredentials([usernamePassword(credentialsId: '21d3d00f-32cd-4d8c-b6ef-2838adf3fb12', usernameVariable: 'REGISTRY_USER', passwordVariable: 'REGISTRY_PASS')]) {
              sh "echo ${REGISTRY_PASS} | docker login -u ${REGISTRY_USER} --password-stdin ${DOCKER_REGISTRY}"
              sh "docker build -t ${IMAGE_FULL_NAME}:${IMAGE_VERSION} ."
              sh "docker push ${IMAGE_FULL_NAME}:${IMAGE_VERSION}"
              sh "docker tag ${IMAGE_FULL_NAME}:${IMAGE_VERSION} ${IMAGE_FULL_NAME}:latest"
              sh "docker push ${IMAGE_FULL_NAME}:latest"
            }
          } catch(err) {
            echo err.getMessage()
          }
        }
      }
    }
  }
}
