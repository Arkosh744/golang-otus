//go:generate easyjson -all

package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		user := User{}
		err = easyjson.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal user: %w", err)
		}
		result = append(result, user)
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat, len(u))

	for _, user := range u {
		mailDomain := strings.SplitN(user.Email, "@", 2)
		if len(mailDomain) != 2 {
			continue
		}

		if strings.HasSuffix(mailDomain[1], domain) {
			subDomain := strings.ToLower(mailDomain[1])
			result[subDomain]++
		}
	}
	return result, nil
}
