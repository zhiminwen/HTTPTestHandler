# Test program for Knative
A simple test HTTP handler for Knative

## Sample Usage
### CPU hogging
  Usage:  http://host:port/cpuHog?period=120&load=0.8   
  Generate 80% of CPU utilization on the system for 120 seconds   

### Memory hogging
  Usage:  http://host:port/memHog?period=120&size=5   
  Allocate 5MB memory and hold a period of 120 seconds

### Time hogging
  Usage: http://host:port/timeHog?period=120   
  Simply hold the HTTP response for 120 seconds
  
### Blue-green deployment
  Set the HTML Body background color based on the Env variable "BGROUND"
