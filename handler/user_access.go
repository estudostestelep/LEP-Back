package handler

import (
	"errors"
	"lep/repositories"
	"lep/repositories/models"
	"time"

	"github.com/google/uuid"
)

type resourceUserAccess struct {
	repo *repositories.DBconn
}

// UserOrganizationAccess representa o acesso de um usuário a uma organização
type UserOrganizationAccess struct {
	Id               uuid.UUID  `json:"id"`
	UserId           uuid.UUID  `json:"user_id"`
	OrganizationId   uuid.UUID  `json:"organization_id"`
	OrganizationName string     `json:"organization_name"`
	Role             string     `json:"role"`
	Active           bool       `json:"active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// UserProjectAccess representa o acesso de um usuário a um projeto
type UserProjectAccess struct {
	Id               uuid.UUID  `json:"id"`
	UserId           uuid.UUID  `json:"user_id"`
	ProjectId        uuid.UUID  `json:"project_id"`
	ProjectName      string     `json:"project_name"`
	OrganizationId   uuid.UUID  `json:"organization_id"`
	OrganizationName string     `json:"organization_name"`
	Role             string     `json:"role"`
	Active           bool       `json:"active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// UserAccessData contém as organizações e projetos que um usuário tem acesso
type UserAccessData struct {
	Organizations []UserOrganizationAccess `json:"organizations"`
	Projects      []UserProjectAccess      `json:"projects"`
}

// UpdateUserAccessRequest representa a requisição para atualizar o acesso de um usuário
type UpdateUserAccessRequest struct {
	OrganizationIds []string `json:"organization_ids"`
	ProjectIds      []string `json:"project_ids"`
}

// UpdateUserAccessResult representa o resultado da atualização de acesso
type UpdateUserAccessResult struct {
	Message              string `json:"message"`
	OrganizationsAdded   int    `json:"organizations_added"`
	OrganizationsRemoved int    `json:"organizations_removed"`
	ProjectsAdded        int    `json:"projects_added"`
	ProjectsRemoved      int    `json:"projects_removed"`
}

type IHandlerUserAccess interface {
	GetUserAccess(userId string) (*UserAccessData, error)
	UpdateUserAccess(userId string, request *UpdateUserAccessRequest) (*UpdateUserAccessResult, error)
	IsAdmin(userId string) (bool, error)
	IsClient(userId string) (bool, error)
}

func (r *resourceUserAccess) IsAdmin(userId string) (bool, error) {
	admin, err := r.repo.Admins.GetAdminById(userId)
	if err != nil {
		return false, nil // Não é admin ou erro ao buscar
	}
	return admin != nil, nil
}

func (r *resourceUserAccess) IsClient(userId string) (bool, error) {
	client, err := r.repo.Clients.GetClientById(userId)
	if err != nil {
		return false, nil // Não é client ou erro ao buscar
	}
	return client != nil, nil
}

func (r *resourceUserAccess) GetUserAccess(userId string) (*UserAccessData, error) {
	result := &UserAccessData{
		Organizations: []UserOrganizationAccess{},
		Projects:      []UserProjectAccess{},
	}

	// Verificar se é admin
	admin, err := r.repo.Admins.GetAdminById(userId)
	if err == nil && admin != nil {
		// É um admin - buscar organizações via AdminRole
		adminRoles, err := r.repo.Roles.GetAdminRoles(userId)
		if err != nil {
			return nil, err
		}

		// Mapear organizações únicas
		orgMap := make(map[uuid.UUID]bool)
		for _, ar := range adminRoles {
			if ar.OrganizationId != nil && !orgMap[*ar.OrganizationId] {
				orgMap[*ar.OrganizationId] = true

				// Buscar detalhes da organização
				org, err := r.repo.Organizations.GetOrganizationById(*ar.OrganizationId)
				if err != nil || org == nil {
					continue
				}

				roleName := ""
				if ar.Role != nil {
					roleName = ar.Role.Name
				}

				result.Organizations = append(result.Organizations, UserOrganizationAccess{
					Id:               ar.Id,
					UserId:           admin.Id,
					OrganizationId:   *ar.OrganizationId,
					OrganizationName: org.Name,
					Role:             roleName,
					Active:           ar.Active,
					CreatedAt:        ar.CreatedAt,
					UpdatedAt:        ar.UpdatedAt,
				})
			}
		}

		// Admin com role global (OrganizationId = null) tem acesso a todas as organizações
		for _, ar := range adminRoles {
			if ar.OrganizationId == nil && ar.Role != nil && ar.Role.HierarchyLevel >= 10 {
				// Super admin - listar todas as organizações
				orgs, err := r.repo.Organizations.ListOrganizations()
				if err == nil {
					for _, org := range orgs {
						if !orgMap[org.Id] {
							orgMap[org.Id] = true
							result.Organizations = append(result.Organizations, UserOrganizationAccess{
								Id:               uuid.New(), // ID virtual
								UserId:           admin.Id,
								OrganizationId:   org.Id,
								OrganizationName: org.Name,
								Role:             "super_admin",
								Active:           true,
								CreatedAt:        org.CreatedAt,
								UpdatedAt:        org.UpdatedAt,
							})
						}
					}
				}
				break
			}
		}

		// Para admins, projetos são acessíveis via organização
		for _, orgAccess := range result.Organizations {
			projects, err := r.repo.Projects.GetProjectByOrganization(orgAccess.OrganizationId)
			if err == nil {
				for _, proj := range projects {
					result.Projects = append(result.Projects, UserProjectAccess{
						Id:               uuid.New(), // ID virtual
						UserId:           admin.Id,
						ProjectId:        proj.Id,
						ProjectName:      proj.Name,
						OrganizationId:   proj.OrganizationId,
						OrganizationName: orgAccess.OrganizationName,
						Role:             orgAccess.Role,
						Active:           true,
						CreatedAt:        proj.CreatedAt,
						UpdatedAt:        proj.UpdatedAt,
					})
				}
			}
		}

		return result, nil
	}

	// Verificar se é client
	client, err := r.repo.Clients.GetClientById(userId)
	if err == nil && client != nil {
		// É um client - organização é fixa (OrgId)
		org, err := r.repo.Organizations.GetOrganizationById(client.OrgId)
		if err == nil && org != nil {
			// Buscar role do client
			clientRoles, _ := r.repo.Roles.GetClientRoles(userId, client.OrgId.String())
			roleName := "member"
			if len(clientRoles) > 0 && clientRoles[0].Role != nil {
				roleName = clientRoles[0].Role.Name
			}

			result.Organizations = append(result.Organizations, UserOrganizationAccess{
				Id:               uuid.New(), // ID virtual
				UserId:           client.Id,
				OrganizationId:   client.OrgId,
				OrganizationName: org.Name,
				Role:             roleName,
				Active:           client.Active,
				CreatedAt:        client.CreatedAt,
				UpdatedAt:        client.UpdatedAt,
			})
		}

		// Projetos do client vêm do ProjIds
		for _, projIdStr := range client.ProjIds {
			projId, err := uuid.Parse(projIdStr)
			if err != nil {
				continue
			}

			proj, err := r.repo.Projects.GetProjectById(projId)
			if err != nil || proj == nil {
				continue
			}

			orgName := ""
			if org != nil {
				orgName = org.Name
			}

			// Buscar role específico do projeto
			clientRoles, _ := r.repo.Roles.GetClientRoles(userId, client.OrgId.String())
			roleName := "member"
			for _, cr := range clientRoles {
				if cr.ProjectId != nil && *cr.ProjectId == projId && cr.Role != nil {
					roleName = cr.Role.Name
					break
				}
			}

			result.Projects = append(result.Projects, UserProjectAccess{
				Id:               uuid.New(), // ID virtual
				UserId:           client.Id,
				ProjectId:        proj.Id,
				ProjectName:      proj.Name,
				OrganizationId:   proj.OrganizationId,
				OrganizationName: orgName,
				Role:             roleName,
				Active:           true,
				CreatedAt:        proj.CreatedAt,
				UpdatedAt:        proj.UpdatedAt,
			})
		}

		return result, nil
	}

	return nil, errors.New("usuário não encontrado")
}

func (r *resourceUserAccess) UpdateUserAccess(userId string, request *UpdateUserAccessRequest) (*UpdateUserAccessResult, error) {
	result := &UpdateUserAccessResult{
		Message: "Acesso atualizado com sucesso",
	}

	// Verificar se é client (atualmente só suportamos atualização de projetos para clients)
	client, err := r.repo.Clients.GetClientById(userId)
	if err != nil || client == nil {
		// Verificar se é admin
		admin, err := r.repo.Admins.GetAdminById(userId)
		if err != nil || admin == nil {
			return nil, errors.New("usuário não encontrado")
		}

		// Para admins, atualizar organizações via AdminRole
		// Buscar roles atuais
		currentRoles, err := r.repo.Roles.GetAdminRoles(userId)
		if err != nil {
			return nil, err
		}

		// Criar mapa de organizações atuais
		currentOrgIds := make(map[string]models.AdminRole)
		for _, role := range currentRoles {
			if role.OrganizationId != nil {
				currentOrgIds[role.OrganizationId.String()] = role
			}
		}

		// Criar mapa de novas organizações
		newOrgIds := make(map[string]bool)
		for _, orgId := range request.OrganizationIds {
			newOrgIds[orgId] = true
		}

		// Remover organizações que não estão mais na lista
		for orgIdStr, role := range currentOrgIds {
			if !newOrgIds[orgIdStr] {
				// Remover role
				if err := r.repo.Roles.RemoveRoleFromAdmin(admin.Id.String(), role.RoleId.String()); err != nil {
					continue
				}
				result.OrganizationsRemoved++
			}
		}

		// Adicionar novas organizações
		// Buscar role padrão para admin (org_admin)
		defaultRole, _ := r.repo.Roles.GetByName("org_admin")
		if defaultRole == nil {
			// Tentar super_admin se org_admin não existir
			defaultRole, _ = r.repo.Roles.GetByName("super_admin")
		}

		for orgIdStr := range newOrgIds {
			if _, exists := currentOrgIds[orgIdStr]; !exists {
				// Adicionar nova associação
				orgId, err := uuid.Parse(orgIdStr)
				if err != nil {
					continue
				}

				roleId := uuid.Nil
				if defaultRole != nil {
					roleId = defaultRole.Id
				}

				if roleId != uuid.Nil {
					adminRole := &models.AdminRole{
						Id:             uuid.New(),
						AdminId:        admin.Id,
						RoleId:         roleId,
						OrganizationId: &orgId,
						Active:         true,
					}
					if err := r.repo.Roles.AssignRoleToAdmin(adminRole); err != nil {
						continue
					}
					result.OrganizationsAdded++
				}
			}
		}

		return result, nil
	}

	// Para clients, atualizar projetos (org é fixa)
	// Validação: usuário deve estar em pelo menos 1 projeto
	if len(request.ProjectIds) == 0 {
		return nil, errors.New("usuário deve estar vinculado a pelo menos 1 projeto")
	}

	// Criar mapa de projetos atuais
	currentProjIds := make(map[string]bool)
	for _, projId := range client.ProjIds {
		currentProjIds[projId] = true
	}

	// Criar mapa de novos projetos
	newProjIds := make(map[string]bool)
	for _, projId := range request.ProjectIds {
		newProjIds[projId] = true
	}

	// Calcular diferenças
	for projId := range currentProjIds {
		if !newProjIds[projId] {
			// Remover projeto
			if err := r.repo.Clients.RemoveProjectFromClient(userId, projId); err == nil {
				result.ProjectsRemoved++
			}
		}
	}

	for projId := range newProjIds {
		if !currentProjIds[projId] {
			// Verificar se projeto pertence à organização do client
			projUUID, err := uuid.Parse(projId)
			if err != nil {
				continue
			}
			proj, err := r.repo.Projects.GetProjectById(projUUID)
			if err != nil || proj == nil || proj.OrganizationId != client.OrgId {
				continue
			}
			// Adicionar projeto
			if err := r.repo.Clients.AddProjectToClient(userId, projId); err == nil {
				result.ProjectsAdded++
			}
		}
	}

	return result, nil
}

func NewUserAccessHandler(repo *repositories.DBconn) IHandlerUserAccess {
	return &resourceUserAccess{repo: repo}
}
