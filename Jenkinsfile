def app
pipeline {

    options {
        disableConcurrentBuilds()
    }

    parameters {
    }

    stages {
        stage('Check') {
            steps {
                withEnv(["PATH+NODEJS=${tool params.NODEJS_TOOL_NAME}/bin"]) {
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
                        sh """
              docker --version
            """

                        BaseimageName = "beeckup/mysql-sidecar"
                    }
                }
            }
        }

        stage('Copy file for docker build') {
            steps {
                sh "cp src/aws.go docker/"
                sh "cp src/common.go docker/"
                sh "cp src/main.go docker/"
                sh "cp src/minio.go docker/"
                sh "cp src/mysql.go docker/"
                sh "cp src/zip.go docker/"
            }
        }
        stage('Build image') {
            steps {
                app = docker.build(BaseimageName, "--pull docker/")
            }
        }

        stage('Push image') {
            steps {
                docker.withRegistry('https://registry.hub.docker.com', 'docker-hub-credentials') {
                    app.push("${tag}")
                }
            }
        }

        stage('Clean') {
            steps {
                sh """
            docker rmi ${BaseimageName}:${tag}
          """
            }
        }
    }

}
