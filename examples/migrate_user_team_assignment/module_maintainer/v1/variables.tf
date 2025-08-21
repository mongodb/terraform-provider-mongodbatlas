variable "org_id" {  
  description = "MongoDB Atlas Organization ID"  
  type        = string  
}  
  
variable "team_name" {  
  description = "Name of the team"  
  type        = string  
}  
  
variable "usernames" {  
  description = "List of usernames to assign to the team"  
  type        = list(string)  
  default     = []  
} 