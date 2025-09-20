# LEP System - GCP Production Environment
# Full production configuration with all features enabled

# Project configuration
project_id   = "leps-472702"
project_name = "leps"
region       = "us-central1"
environment  = "prod"

# Database configuration - FULL PRODUCTION
database_name              = "lep_database_prod"
database_user              = "lep_user_prod"
db_tier                    = "db-n1-standard-2"  # High-performance tier
db_availability_type       = "REGIONAL"          # High availability
db_disk_size               = 100                  # Large disk for production data
enable_deletion_protection = true                # Protect against accidental deletion

# JWT configuration
jwt_private_key = "-----BEGIN ENCRYPTED PRIVATE KEY-----\nMIIFNTBfBgkqhkiG9w0BBQ0wUjAxBgkqhkiG9w0BBQwwJAQQznM3L7SM4uu1WbgK\naqBgLQICCAAwDAYIKoZIhvcNAgkFADAdBglghkgBZQMEASoEEKGmGvATnEoRiVa7\n2zH7HSEEggTQjqqPViJBswLieglxj5HxdVNN/9qd+EgWRoF0xYgfesStW/ahQ/ly\nc16pVIq1IAvuo3BmAVVOnWClBet1TE28QOV8qbfe24msj0vLvCb7xhta3eyXLj2A\niPEa/d5h7updctxwQ9n4YP0a7PkgcgL7BNgJrmEtLuT9JuytTXzNLVjheOQj+SCa\npO44ScJZY1aMTyLy1FT3q/zAIpw/knZbf0gBolSGso4ERrvVirgOLrUS/GnfwhWp\nx+OVY4rsvSZNUHQ7sbMBqf//R4b5DhInwTMYb4Pao0dsh4PCCUKYtTp0SMakBInY\n0f+LYtr+TZCI3evGtyNpU0gsfA1v3rG2VCe9dQgrB7miGG7SVMuqIht4SFL1145A\nGaXECXQEgKukhVZbiblFcEHSp1MGEUxcbWWUyM4dQ0jTB3v+Kh4KQmwbBoU4tnAs\nmsf4/+f1rBsGtX/HW5XvZv2tf/bZW8hiklD0mEEn+3ta4FA5y+131zCViE/TknPV\nke2jcyfgg3HcnR2f+j150CUf5pJeom+FpF9ngBajrMG/sukggLukCrtYilFNiFiR\nfgcQ4PG2wjsfAphcOTmR+LhCosSzo+HzDPQ7+gyjKvoXecZIVJuw4DKCJEuvBVrg\n+i2jFceft4mfCobwmpLnWifzLn9RcLCQ079TjWDLNnvNrXcgZ3A/3qu8Np07Oa6b\nS6Qcl+3/JVHtdROHMnNrHwUAlWwKVVEtrgsuMzXXuA2WM6firOWsxC7ozM7z53zT\nUJ31BeCeaskoDaNFocMbZXoYxMoVEgrOOlJLYBWKa9Lg1VRCgP6Uj3Ob0OiNx2Nz\nkXsEoIS1iN+31UEIipWqP9bqNBwfSraaM8mxCMsbdqRZAhj4/Qeup2nW4IRoWMCM\nRP643D97WqJsl9VPATK+ujwO/GVpjJdxoAMLn2jEqGnubat7AMeibHPwPqoqikFb\nCnflm4XnQJO475ENqyTQ4QB2OfWwuAJaIDbCde1gTvcthqd0SJtzGySi21kHb11n\nRFIdH1SSqR0H6qcRgeEnKeXC4ANNtCkhZSCaGnSBRvrBb1yfugogND/UHf8hV38y\nOCFn0YL/Umur61jmWxNns8TdFtZoT44LfCzd2i5aIs7IWmPc63Skp9E1fMZswsUC\nXp2nzHAS5jvtWfGD87w8NrmnRqPfXlApX0nBDHmy2HyF5/IWkmyhQ3PPzc6hlZgb\nt+vTov8RwE15aLEPRyY/p564sTkg8xiYIumL9KUT0aWLcxFG7uViH/CjZ3hmrj7t\nbi5FYoBQ+MAr4HJ5wKWn+xtkc7EuUPxl+biCcpjnb4KuwpHgGQVRv/Fa6drt2ma6\nb01w8Dbu+FLGqkYpSrBNTga/93vpSOYim3Yd4sJlLrDYjrrrwRl4mkFHHjDSspZz\nLggCSKKcpZJ2uaEhmNfdC+KVgGLILdUPPtVgm2XM6QAnJs/PjJ6PZN5OZIWSr3dO\nRKunUlpZyM3KULhxljeeP36rrzYBxwhz1mQhI+IBPBsT3A9IkXAe3OyO2+58eyaG\nYHw5fFsZY1LLiopwbiMo6i3dFNh2saeFqu7jpH1Ag6+wbN8D83G9+qZ0SQ61tLDW\n7fHVeZOJ0GFQmUoc7sUMGr7iTxSokSbywh2wce8n6j7oKCyYJg2AOrE=\n-----END ENCRYPTED PRIVATE KEY-----"
jwt_public_key  = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsyilGx49Drplmoq8Uqfr\ndNa5xE5pfL9eAaTirL+5Ah0RWucUmZHh5NJpY3FpYS6GlxQh56KiQf6PsqOzEarM\nKBkwGytEGjER7osaC7n6f9aZduwtlmhibJTuthYr1WTcV/cR8E5Lv81wfeLSYTWm\nK3Ek0P3sv1MWL52Khthn9Sdus3YRH5Nocbo+TarleHa68zQY8DZZzGFIHqB7Hl96\nIoErLp+OMbqxnYy111G+XxRkdxg0GDsxgl7DHZj/4Ucz8PJ58owAlYsMJNNfPzGQ\nbf3Op3xtlEeKptCdv57MemKu+gOjaFjfDmE5opljaKD2zeNMhw7IuDaFEBjBq6ma\nMQIDAQAB\n-----END PUBLIC KEY-----"

# Cloud Run configuration - FULL PRODUCTION
min_instances = 2          # Always keep instances warm
max_instances = 50         # High capacity for production load
cpu_limit     = "4"        # Maximum CPU power
memory_limit  = "2Gi"      # Maximum memory

# Twilio configuration - ENABLED for full production
twilio_account_sid  = "your-twilio-account-sid"      # UPDATE THIS
twilio_auth_token   = "your-twilio-auth-token"       # UPDATE THIS
twilio_phone_number = "+1234567890"                  # UPDATE THIS

# SMTP configuration - ENABLED for production emails
smtp_host     = "smtp.gmail.com"
smtp_port     = 587
smtp_username = "your-production-email@gmail.com"    # UPDATE THIS
smtp_password = "your-app-password"                  # UPDATE THIS

# Application configuration
enable_cron_jobs = true    # Enable all features

# Custom domain - ENABLED for production
domain_name          = "api.yourdomain.com"         # UPDATE THIS
enable_custom_domain = true                         # Enable custom domain

# Additional production settings
# Note: These would be additional variables if needed
# monitoring_enabled = true
# backup_retention_days = 30
# ssl_redirect = true