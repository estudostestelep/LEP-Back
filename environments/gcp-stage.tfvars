# LEP System - GCP Staging Environment
# Production-like configuration without Twilio (SMTP only)

# Project configuration
project_id   = "leps-472702"
project_name = "leps"
region       = "us-central1"
environment  = "staging"

# Database configuration - PRODUCTION-LIKE
database_name              = "lep_database_staging"
database_user              = "lep_user_staging"
db_tier                    = "db-n1-standard-1"  # Production-ready tier
db_availability_type       = "REGIONAL"          # High availability
db_disk_size               = 50                   # Larger disk
enable_deletion_protection = true                # Protect against accidental deletion

# JWT configuration
jwt_private_key = "-----BEGIN ENCRYPTED PRIVATE KEY-----\nMIIFNTBfBgkqhkiG9w0BBQ0wUjAxBgkqhkiG9w0BBQwwJAQQznM3L7SM4uu1WbgK\naqBgLQICCAAwDAYIKoZIhvcNAgkFADAdBglghkgBZQMEASoEEKGmGvATnEoRiVa7\n2zH7HSEEggTQjqqPViJBswLieglxj5HxdVNN/9qd+EgWRoF0xYgfesStW/ahQ/ly\nc16pVIq1IAvuo3BmAVVOnWClBet1TE28QOV8qbfe24msj0vLvCb7xhta3eyXLj2A\niPEa/d5h7updctxwQ9n4YP0a7PkgcgL7BNgJrmEtLuT9JuytTXzNLVjheOQj+SCa\npO44ScJZY1aMTyLy1FT3q/zAIpw/knZbf0gBolSGso4ERrvVirgOLrUS/GnfwhWp\nx+OVY4rsvSZNUHQ7sbMBqf//R4b5DhInwTMYb4Pao0dsh4PCCUKYtTp0SMakBInY\n0f+LYtr+TZCI3evGtyNpU0gsfA1v3rG2VCe9dQgrB7miGG7SVMuqIht4SFL1145A\nGaXECXQEgKukhVZbiblFcEHSp1MGEUxcbWWUyM4dQ0jTB3v+Kh4KQmwbBoU4tnAs\nmsf4/+f1rBsGtX/HW5XvZv2tf/bZW8hiklD0mEEn+3ta4FA5y+131zCViE/TknPV\nke2jcyfgg3HcnR2f+j150CUf5pJeom+FpF9ngBajrMG/sukggLukCrtYilFNiFiR\nfgcQ4PG2wjsfAphcOTmR+LhCosSzo+HzDPQ7+gyjKvoXecZIVJuw4DKCJEuvBVrg\n+i2jFceft4mfCobwmpLnWifzLn9RcLCQ079TjWDLNnvNrXcgZ3A/3qu8Np07Oa6b\nS6Qcl+3/JVHtdROHMnNrHwUAlWwKVVEtrgsuMzXXuA2WM6firOWsxC7ozM7z53zT\nUJ31BeCeaskoDaNFocMbZXoYxMoVEgrOOlJLYBWKa9Lg1VRCgP6Uj3Ob0OiNx2Nz\nkXsEoIS1iN+31UEIipWqP9bqNBwfSraaM8mxCMsbdqRZAhj4/Qeup2nW4IRoWMCM\nRP643D97WqJsl9VPATK+ujwO/GVpjJdxoAMLn2jEqGnubat7AMeibHPwPqoqikFb\nCnflm4XnQJO475ENqyTQ4QB2OfWwuAJaIDbCde1gTvcthqd0SJtzGySi21kHb11n\nRFIdH1SSqR0H6qcRgeEnKeXC4ANNtCkhZSCaGnSBRvrBb1yfugogND/UHf8hV38y\nOCFn0YL/Umur61jmWxNns8TdFtZoT44LfCzd2i5aIs7IWmPc63Skp9E1fMZswsUC\nXp2nzHAS5jvtWfGD87w8NrmnRqPfXlApX0nBDHmy2HyF5/IWkmyhQ3PPzc6hlZgb\nt+vTov8RwE15aLEPRyY/p564sTkg8xiYIumL9KUT0aWLcxFG7uViH/CjZ3hmrj7t\nbi5FYoBQ+MAr4HJ5wKWn+xtkc7EuUPxl+biCcpjnb4KuwpHgGQVRv/Fa6drt2ma6\nb01w8Dbu+FLGqkYpSrBNTga/93vpSOYim3Yd4sJlLrDYjrrrwRl4mkFHHjDSspZz\nLggCSKKcpZJ2uaEhmNfdC+KVgGLILdUPPtVgm2XM6QAnJs/PjJ6PZN5OZIWSr3dO\nRKunUlpZyM3KULhxljeeP36rrzYBxwhz1mQhI+IBPBsT3A9IkXAe3OyO2+58eyaG\nYHw5fFsZY1LLiopwbiMo6i3dFNh2saeFqu7jpH1Ag6+wbN8D83G9+qZ0SQ61tLDW\n7fHVeZOJ0GFQmUoc7sUMGr7iTxSokSbywh2wce8n6j7oKCyYJg2AOrE=\n-----END ENCRYPTED PRIVATE KEY-----"
jwt_public_key  = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsyilGx49Drplmoq8Uqfr\ndNa5xE5pfL9eAaTirL+5Ah0RWucUmZHh5NJpY3FpYS6GlxQh56KiQf6PsqOzEarM\nKBkwGytEGjER7osaC7n6f9aZduwtlmhibJTuthYr1WTcV/cR8E5Lv81wfeLSYTWm\nK3Ek0P3sv1MWL52Khthn9Sdus3YRH5Nocbo+TarleHa68zQY8DZZzGFIHqB7Hl96\nIoErLp+OMbqxnYy111G+XxRkdxg0GDsxgl7DHZj/4Ucz8PJ58owAlYsMJNNfPzGQ\nbf3Op3xtlEeKptCdv57MemKu+gOjaFjfDmE5opljaKD2zeNMhw7IuDaFEBjBq6ma\nMQIDAQAB\n-----END PUBLIC KEY-----"

# Cloud Run configuration - PRODUCTION-LIKE
min_instances = 1          # Always keep one instance warm
max_instances = 20         # Higher capacity for load testing
cpu_limit     = "2"        # More CPU power
memory_limit  = "1Gi"      # More memory

# Twilio configuration - DISABLED for staging
twilio_account_sid  = ""
twilio_auth_token   = ""
twilio_phone_number = ""

# SMTP configuration - ENABLED for email testing
smtp_host     = "smtp.gmail.com"
smtp_port     = 587
smtp_username = "your-staging-email@gmail.com"    # UPDATE THIS
smtp_password = "your-app-password"                # UPDATE THIS

# Application configuration
enable_cron_jobs = true    # Enable for staging testing

# Custom domain - optional for staging
domain_name          = "staging-api.yourdomain.com"  # UPDATE THIS if needed
enable_custom_domain = false                         # Set to true if domain is configured