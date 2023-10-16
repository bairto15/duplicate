package service


func (s *Service) IsDuplicate(userId1, userId2 int32) bool {
	s.logger.Info(userId1)
	s.logger.Info(userId2)

	if userId1 == userId2 {
		return true
	}

	s.mutex.RLock()

	cash := s.repo.GetCash()
	setIps := cash.UsersWithIps[userId1] //Получаем список ip адресов по данному user_id

	for ipAddr := range setIps {
		setUserIds, ok := cash.IpsWithUsers[ipAddr] //Получаем список user_id по ip адресу
		if ok {
			_, ok := setUserIds[userId2] //Проверяем есть ли в списке user_id
			if ok {
				s.logger.Info(ipAddr)
				return true
			}
		}
	}

	s.mutex.RUnlock()

	return false
}
