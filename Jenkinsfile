pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        echo 'Building...'
        sh 'go build -v'
      }
    }
    stage('Test') {
      steps {
        echo 'Testing...'
        sh 'go test -v'
      }
    }
  }
}
