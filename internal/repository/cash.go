package repository


// Сохранить лог в кэш для быстрого доступа
func (p *Postgres) SaveLogInCash() error {
	users := make(map[int32]SetIps, 1000)
	ips := make(map[string]SetUserId, 1000)

	query := "SELECT user_id, ip_addr FROM conn_log"
	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Error(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var row log

		err = rows.Scan(&row.UserId, &row.IpAddr)
		if err != nil {
			p.logger.Error(err)
			continue
		}

		//Сохранить в мапе по ключу user_id по значению сет из ip адресов
		setIps, ok := users[row.UserId]
		if ok {
			setIps[row.IpAddr] = struct{}{}
			users[row.UserId] = setIps
		} else {
			addr := make(map[string]struct{})
			addr[row.IpAddr] = struct{}{}
			users[row.UserId] = addr
		}

		//Сохранить в мапе по ключу ip адрес по значению 
		setUsers, ok := ips[row.IpAddr]
		if ok {
			setUsers[row.UserId] = struct{}{}
			ips[row.IpAddr] = setUsers
		} else {
			id := make(map[int32]struct{})
			id[row.UserId] = struct{}{}
			ips[row.IpAddr] = id
		}
	}

	p.Cash.UsersWithIps = users
	p.Cash.IpsWithUsers = ips

	return nil
}

func (p *Postgres) GetCash() Cash {
	return p.Cash
}