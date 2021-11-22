pipeline {
    agent {
        docker { 
            // TODO: use a custom docker image in my own aws registry
            image 'takasago/kaizen-jenkins-builder:latest' 
            args '-u root --privileged -v /var/run/docker.sock:/var/run/docker.sock'
        }
    }

    environment {
        COMMIT="${GIT_COMMIT}"
        TAG="${TAG_NAME}"
    }

    stages {
        stage('install'){
            steps {
                sh '''
                make lint-install
                '''
            }
        }

        stage ('checks'){
            failFast true
            parallel {
                stage('tests'){
                    steps {
                        sh '''
                            make test
                        '''
                    }
                }
                stage('linter'){
                    steps {
                        sh '''
                            make lint
                        '''
                    }
                }
            }
        }
    }
}