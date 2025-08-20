variable "org_id" {  
  type        = string  
  description = "MongoDB Atlas Organization ID"  
}  
  
variable "team_name" {  
  type        = string  
  description = "Name of the Atlas team"  
}  
  
variable "usernames" {  
  type        = list(string)  
  description = "List of user emails to assign to the team"  
}  