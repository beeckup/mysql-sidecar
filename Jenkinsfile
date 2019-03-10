def app
def tag
pipeline {

    options {
        disableConcurrentBuilds()
    }

    agent any

    stages {
        stage('Check') {
            steps {

                    script {
                        // enforce branches
                        switch (BRANCH_NAME) {
                            case "master":
                                tag = "latest"
                                break
                            case "test":
                                tag = "test"
                                break
                            case "dev":
                                tag = "dev"
                                break
                            default:
                                error("Error")
                                break
                        }

                        // check tools
                        BaseimageName = "beeckup/mysql-sidecar"
                    }

            }
        }

        stage('Copy file for docker build') {
            steps {
                sh "cp -r src docker/src"
            }
        }
        stage('Build image') {
            steps {
                script {
                    app = docker.build(BaseimageName, "--pull docker/")
                }
            }
        }

        stage('Push image') {
            steps {
                script {
                    docker.withRegistry('https://registry.hub.docker.com', 'docker-hub-credentials') {
                        app.push("${tag}")
                        app.push("${env.BUILD_NUMBER}")
                    }
                }
            }
        }

        stage('Clean') {
            steps {
                sh """
            docker rmi registry.hub.docker.com/${BaseimageName}:${tag}
            docker rmi registry.hub.docker.com/${BaseimageName}:${env.BUILD_NUMBER}
          """
            }
        }
    }

}
