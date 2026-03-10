package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// ResponseClassifierService classifica respostas de clientes em português brasileiro
type ResponseClassifierService struct{}

// ClassificationResult resultado da classificação de uma resposta
type ClassificationResult struct {
	ResponseType    string  // "confirmed", "cancelled", "unknown"
	ConfidenceScore float64 // 0.0 a 1.0
	Method          string  // "pattern_match", "ai_classification"
}

// Padrões de alta confiança para confirmação (exact match)
var confirmPatternsExact = []string{
	`^sim$`, `^s$`, `^ss$`, `^sss+$`,
	`^ok$`, `^okay$`, `^blz$`, `^beleza$`,
	`^confirmo$`, `^confirmado$`, `^confirmada$`,
	`^pode confirmar$`, `^confirma$`,
	`^vou$`, `^vou sim$`, `^vou la$`,
	`^irei$`, `^estarei la$`, `^estarei presente$`,
	`^comparecerei$`, `^vou comparecer$`,
	`^pode ser$`, `^pode$`, `^bora$`,
	`^fechado$`, `^combinado$`, `^certo$`,
	`^com certeza$`, `^claro$`, `^obvio$`,
	`^positivo$`, `^afirmativo$`,
	`^ta bom$`, `^tudo bem$`, `^tudo certo$`,
	`^perfeito$`, `^otimo$`, `^maravilha$`,
	`^👍$`, `^✅$`, `^✔️$`, `^✔$`,
}

// Padrões de média confiança para confirmação (partial match)
var confirmPatternsPartial = []string{
	`confirm`, `vou.*sim`, `estarei`, `irei.*sim`,
}

// Padrões de alta confiança para cancelamento (exact match)
var cancelPatternsExact = []string{
	`^nao$`, `^n$`, `^nn$`, `^nnn+$`,
	`^cancela$`, `^cancelar$`, `^cancele$`,
	`^cancelado$`, `^cancelada$`,
	`^nao vou$`, `^nao irei$`, `^nao posso$`,
	`^nao vou poder$`, `^nao da$`, `^nao vai dar$`,
	`^desisto$`, `^desistir$`,
	`^nao quero$`, `^nao quero mais$`,
	`^nao comparecerei$`,
	`^tive um imprevisto$`, `^imprevisto$`,
	`^surgiu algo$`, `^nao vai rolar$`,
	`^fica pra proxima$`, `^outra vez$`,
	`^negativo$`,
	`^👎$`, `^❌$`, `^✖️$`, `^✖$`,
}

// Padrões de média confiança para cancelamento (partial match)
var cancelPatternsPartial = []string{
	`cancel`, `nao.*poder`, `desist`, `nao.*ir`, `imprevisto`,
}

// NewResponseClassifierService cria nova instância do classificador
func NewResponseClassifierService() *ResponseClassifierService {
	return &ResponseClassifierService{}
}

// ClassifyResponse analisa a mensagem e retorna a classificação
func (r *ResponseClassifierService) ClassifyResponse(message string) ClassificationResult {
	// Normaliza a mensagem: minúsculas, remove acentos, trim
	normalized := r.normalizeText(message)

	// Tenta match exato primeiro (alta confiança)
	if r.matchPatterns(normalized, confirmPatternsExact) {
		return ClassificationResult{
			ResponseType:    "confirmed",
			ConfidenceScore: 0.95,
			Method:          "pattern_match",
		}
	}

	if r.matchPatterns(normalized, cancelPatternsExact) {
		return ClassificationResult{
			ResponseType:    "cancelled",
			ConfidenceScore: 0.95,
			Method:          "pattern_match",
		}
	}

	// Tenta match parcial (média confiança)
	if r.matchPatterns(normalized, confirmPatternsPartial) {
		return ClassificationResult{
			ResponseType:    "confirmed",
			ConfidenceScore: 0.75,
			Method:          "pattern_match",
		}
	}

	if r.matchPatterns(normalized, cancelPatternsPartial) {
		return ClassificationResult{
			ResponseType:    "cancelled",
			ConfidenceScore: 0.75,
			Method:          "pattern_match",
		}
	}

	// Nenhum match encontrado
	return ClassificationResult{
		ResponseType:    "unknown",
		ConfidenceScore: 0.0,
		Method:          "pattern_match",
	}
}

// normalizeText normaliza o texto removendo acentos e convertendo para minúsculas
func (r *ResponseClassifierService) normalizeText(text string) string {
	// Minúsculas
	text = strings.ToLower(text)

	// Remove acentos
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, text)

	// Trim espaços
	result = strings.TrimSpace(result)

	return result
}

// matchPatterns verifica se o texto corresponde a algum dos padrões
func (r *ResponseClassifierService) matchPatterns(text string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, text)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// ClassifyWithAI classificação usando IA (placeholder para expansão futura)
// Pode ser expandido para chamar API do Claude ou outro serviço de IA
func (r *ResponseClassifierService) ClassifyWithAI(message string, context map[string]string) ClassificationResult {
	// Por enquanto, retorna unknown - pode ser expandido para integração com IA
	return ClassificationResult{
		ResponseType:    "unknown",
		ConfidenceScore: 0.0,
		Method:          "ai_classification",
	}
}
