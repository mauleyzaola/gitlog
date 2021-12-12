def isProd

if (BRANCH_NAME =~ /^(master)$/)  {
    isProd = true
} else{
    isProd = false
}

pipeline {
    agent {
        docker { 
            // TODO: use a custom docker image in my own aws registry
            image 'mauleyzaola/docker-golang:latest'
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
                    when {
                        expression {
                            return !isProd
                        }
                    }
                    steps {
                        sh '''
                            make test
                        '''
                    }
                }
                stage('linter'){
                    when {
                        expression {
                            return !isProd
                        }
                    }
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