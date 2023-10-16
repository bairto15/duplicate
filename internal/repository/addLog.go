package repository


//Записать новый лог в базу и добавить в кэш
func (p *Postgres) AddLog(log log) error {
	query := "INSERT INTO conn_log (user_id, ip_addr, ts) VALUES ($1, $2, 'NOW()')"
	err := p.db.QueryRow(query, log.UserId, log.IpAddr).Err()
	if err != nil {
		p.logger.Error(err)
		return err
	}

	users := p.Cash.UsersWithIps
	ips := p.Cash.IpsWithUsers

	p.mutex.Lock()

	setIps, ok := users[log.UserId]
	if ok {
		setIps[log.IpAddr] = struct{}{}
		users[log.UserId] = setIps
	} else {
		addr := make(map[string]struct{})
		addr[log.IpAddr] = struct{}{}
		users[log.UserId] = addr
	}

	//Сохранить в мапе по ключу ip адрес по значению
	setUsers, ok := ips[log.IpAddr]
	if ok {
		setUsers[log.UserId] = struct{}{}
		ips[log.IpAddr] = setUsers
	} else {
		id := make(map[int32]struct{})
		id[log.UserId] = struct{}{}
		ips[log.IpAddr] = id
	}

	p.Cash.UsersWithIps = users
	p.Cash.IpsWithUsers = ips

	p.mutex.Unlock()

	return nil
}
