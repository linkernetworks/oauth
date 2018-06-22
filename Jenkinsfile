pipeline {
    agent {
        dockerfile {
            dir "src/github.com/linkernetworks/oauth/jenkins"
            args "--privileged --group-add docker"
        }
    }
    post {
        always {
            dir ("src/github.com/linkernetworks/oauth") {
                sh "make clean"
            }
        }
        success {
            script {
                def message =   "<https://jenkins.linkernetworks.co/job/oauth/|oauth> » " +
                                "<${env.JOB_URL}|${env.BRANCH_NAME}> » " +
                                "<${env.BUILD_URL}|#${env.BUILD_NUMBER}> passed."

                slackSend channel: '#09_jenkins', color: 'good', message: message
            }
        }
        failure {
            script {
                def message =   "<https://jenkins.linkernetworks.co/job/oauth/|oauth> » " +
                                "<${env.JOB_URL}|${env.BRANCH_NAME}> » " +
                                "<${env.BUILD_URL}|#${env.BUILD_NUMBER}> failed."

                slackSend channel: '#09_jenkins', color: 'danger', message: message
            }
        }
        fixed {
            slackSend channel: '#09_jenkins', color: 'good',
                message:    "<https://jenkins.linkernetworks.co/job/oauth/|oauth> » " +
                            "<${env.JOB_URL}|${env.BRANCH_NAME}> » " +
                            "<${env.BUILD_URL}|#${env.BUILD_NUMBER}> is fixed."
        }
    }
    options {
        timestamps()
        timeout(time: 1, unit: 'HOURS')
        checkoutToSubdirectory('src/github.com/linkernetworks/oauth')
    }
    stages {
        stage('Prepare') {
            steps {
                withEnv(["GOPATH+AA=${env.WORKSPACE}"]) {
                    dir ("src/github.com/linkernetworks/oauth") {
                        sh "make pre-build"
                    }
                }
            }
        }
        stage('Build') {
            steps {
                withEnv(["GOPATH+AA=${env.WORKSPACE}"]) {
                    dir ("src/github.com/linkernetworks/oauth") {
                        sh "make build"
                    }
                }
            }
        }
        stage('Test') {
            steps {
                withEnv(["GOPATH+AA=${env.WORKSPACE}"]) {
                    dir ("src/github.com/linkernetworks/oauth") {
                        sh "make src.test-coverage 2>&1 | tee >(go-junit-report > report.xml)"
                        junit "report.xml"
                        sh 'gocover-cobertura < build/src/coverage.txt > cobertura.xml'
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
        }
        stage("Build Image"){
            steps {
                script {
                    dir ("src/github.com/linkernetworks/oauth") {
                        docker.build("linkernetworks/oauth:${env.BRANCH_NAME.replaceAll("[^A-Za-z0-9.]", "-").toLowerCase()}")
                    }
                }
            }
        }
        stage("Push Image"){
            when {
                branch 'master'
            }
            steps {
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
                    docker.image("linkernetworks/oauth:${env.BRANCH_NAME.replaceAll("[^A-Za-z0-9.]", "-").toLowerCase()}").push("latest")
                }
            }
        }
    }
}