def slack(color, message) {
    echo "${message}"
    slackSend channel: '#09_jenkins',
        color: color,
        message:
            "<${JOB_DISPLAY_URL}|oauth> » " +
            "<${env.JOB_URL}|${env.BRANCH_NAME}> » " +
            "<${env.RUN_DISPLAY_URL}|${env.BUILD_DISPLAY_NAME}> ${message}"
}

def shouldDeploy () {
    switch (env.BRANCH_NAME) {
        case ~/(.+-)?rc(-.+)?$/:
        case 'develop':
        case 'master':
            return true
        default:
            return false
    }
}

pipeline {
    agent {
        kubernetes {
            cloud 'kubernetes'
            // this guarantees the agent will use this template
            label "oauth-pod-${UUID.randomUUID()}"
            defaultContainer 'golang'
            yaml """
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: golang
    image: golang
    imagePullPolicy: Always
    command:
    - cat
    tty: true
  - name: docker
    image: docker:dind
    imagePullPolicy: Always
    securityContext:
      privileged: true
    ports:
    - containerPort: 2375
    tty: true
    command: [ "sh", "-c", "apk add bash && dockerd-entrypoint.sh" ]
"""
        }
    }

    options {
        timestamps()
        timeout(time: 1, unit: 'HOURS')
    }
    stages {
        stage('Wait service') {
            failFast false
            parallel {
                stage('docker') {
                    steps {
                        container('docker') {
                            script {
                                waitUntil {
                                    0 == sh(script:"docker run --rm hello-world", returnStatus: true)
                                }
                            }
                        }
                    }
                }
            }
        }
        stage('Prepare') {
            failFast false
            parallel {
                stage('govendor') {
                    steps {
                        container('golang') {
                            sh "go get -u github.com/kardianos/govendor"
                            sh "go get -u github.com/jstemmer/go-junit-report"
                            sh "go get -u github.com/t-yuki/gocover-cobertura"
                            sh """
                                mkdir -p /go/src/github.com/linkernetworks
                                ln -s `pwd` /go/src/github.com/linkernetworks/oauth
                                cd /go/src/github.com/linkernetworks/oauth
                                make pre-build
                            """
                        }
                    }
                }
            }
        }
        stage('Build') {
            failFast false
            parallel {
                stage('golang') {
                    steps {
                        container('golang') {
                            sh """
                                cd /go/src/github.com/linkernetworks/oauth
                                make build
                            """
                        }
                    }
                }
                stage('docker image') {
                    steps {
                        container('docker') {
                            script {
                                docker.build("linkernetworks/oauth:${env.BRANCH_NAME.replaceAll("[^A-Za-z0-9.]", "-").toLowerCase()}")
                            }
                        }
                    }
                }
            }
        }
        stage('Test') {
            failFast false
            parallel {
                stage ("golang test") {
                    steps {
                        container('golang') {

                            sh """
                                cd /go/src/github.com/linkernetworks/oauth
                                make src.test-coverage 2>&1 | tee >(go-junit-report > report.xml)
                            """
                            junit "report.xml"

                            sh """
                                cd /go/src/github.com/linkernetworks/oauth
                                gocover-cobertura < build/src/coverage.txt > cobertura.xml
                            """

                            cobertura coberturaReportFile: "cobertura.xml", failNoReports: true, failUnstable: true
                            publishHTML (target: [
                                allowMissing: true,
                                alwaysLinkToLastBuild: true,
                                keepAll: true,
                                reportDir: 'build/src',
                                reportFiles: 'coverage.html',
                                reportName: "GO cover report",
                                reportTitles: "GO cover report",
                                includes: "coverage.html"
                            ])
                        }
                    }
                }
                stage ("bats test") {
                    steps {
                        echo "TODO: execute bats test"
                    }
                }
            }
        }
        stage('Deploy') {
            when {
                expression { -> shouldDeploy() }
            }
            stages {
                stage("Push Image"){
                    steps {
                        container('docker') {
                            script {
                                withCredentials([
                                    usernamePassword(
                                        credentialsId: 'docker_hub_linkernetworks',
                                        usernameVariable: 'DOCKER_USER',
                                        passwordVariable: 'DOCKER_PASS'
                                    )
                                ]) {
                                    sh 'echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin'
                                }
                                docker.image(
                                    "linkernetworks/oauth:${env.BRANCH_NAME.replaceAll("[^A-Za-z0-9.]", "-").toLowerCase()}"
                                ).push(
                                    "${env.BRANCH_NAME.replaceAll("[^A-Za-z0-9.]", "-").toLowerCase()}-latest"
                                )
                            }
                        }
                    }
                }
                stage("Deploy"){
                    steps {
                        build job: "oauth/bitbucket-oauth/${env.BRANCH_NAME}", wait: false
                    }
                }
            }
        }
    }
    post{
        success {
            slack('good', "Successed")
        }
        failure {
            slack('danger', "Failed")
        }
        aborted {
            slack('warning', "Aborted")
        }
    }
}
