package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// CONFIGURATION
// ============================================================================

const (
	GraphQLEndpoint = "https://bff-page.kakao.com/graphql"

	// Thay Cookie mới của bạn vào đây
	CookieString = "_kpdid=193cdcc72d6d4125917f8fa486d1337f; _kpiid=26fd12766e8056921aa387ec10f166af; _kau=b875a0dd52f9098dc3fd08a5c6fa113ad0c2542a2da5068b9bf75a5cc6375c547dc9436f99d8c93ff1a36dfbb2d3038b897e2918dfd23973f18e276dd89a0c6ea6e7a2e3f071a7b77d6b4e676f4e27c50ddc92e0ee2105d5a810607feffd631d69879f4b18836105acc22a5d74575ac18bb09428ea58719308b147df37373530363037383639313238313538333336343930343536383835373735391fe573b1bed976a0f02ed351416491cb; _kawlt=ht44UyEO1BUNUhCDNVMtg0R5RdjJ9UOvLkw2IMbWrt7JQxpbCN7qbyHn3Su2YgB658JuB6r_YkVyUf82HzlfVbvAPIfWY6uu8StZMpx1hoXJMWaZ2HVfNfuMOBcIig6b; _kawltea=1764846374; _karmt=K6PdP8gOK9toBl1uDvv6-X5gwpDcCinlxs4y-i8SOxhr9hxhTsqyPjhPMGB9lKmM; _karmtea=1767373574; _kahai=ba4be1e91c915262b60a6402f5e1e321446b1a1739f6a79ad21e84f1393c8927; _kaslt=C+vNz6HH4bazZZ54gPvRZVMCtbWnOjnDglkP05OtWo7TtmhPGhSB7KmPE3KcRG5K0CNE2Ij/en2IfcAoBrSEjA==; _kpawbat_e=BPhCIMKtzdQYfN0IBU84X0Vw0HHAcnQsbjOBC0ItthkofB8vrJhHOYm4Qo8JTmjyNPKDhskGzhZDne1RG1JHauIM3NysoabtJdCa%2B9UNSo0%3D; _kpwtkn=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..ysC1G0cP1czkN1BF.RLUWNMFtFLh7-qdZ_9735-9o0ZBNQJulkDhxb0wG7xxRG7efCFagsz87Xj8GiH-a-Y6J_LR2MrdO0Qz_bWbcBEAkESg9nLjDOqNsP2aJZVr2WVutnUJCSF04gpSPEWbYBV6TUCsRxqSqCBjgdNpWRS97zu81twFbBRD3FVfMmGYN14VvFNiQdDQlK_UGiesa9HIFu-iSuuOKW3U-z9Yj1zalKVvT0YbTAQ33oGNcOomsDxow9GOxtOFx6KqPvcleC98PJ4THJ9lXQyfh_IlX9d5ixNfLit6QP5XhFq35PL20BoJ7VsykRKR3zsvMhmfAszJFFiiYzDVwT7oZ9g9QPwGINS8uoGJSpV5yHa3qOtA6mvF-9qU3GalqCSO2OROiyw1O49-wN4h5l-EMtpC3xM8tajrQTwdeyVtR3TM57hRnZ-jglJ3f4-o34J_jkUHIOpEINaUXUU2F6lTCbOM2DfudBeln7zMJz5sp8ulmuVDB1esWMK5SE_QupuHP8_32QPawPSJifkbn77J3I2NHwxM6VrgnKzV46sb_HNy9H8eYvtsf88qAW4wvpLeK50EKR96Wv7EOZT-gH8yWlXKt4iElGQOA7fj0NGEDronNmOMewk7rSS4f6YYJoqHSI2tvDIQXsYEDnKhhuejbKWBskbhkPY4vAq1ZxwXQKpOfEmvFT3rc6w2P78SLyFqtgnupR-VdKqrhONawZxPpfVvqwyaGpyxn7t-JR2LzJ-gfrNGpsF79mkz75zmyaHWl4W7ufb_IxQtEz81mFwr2HhWyOjFDhCfHy9LqMvNLNxCZxIUXriMJ8aFQziqIbHQYhcSxtnbvtNLkB3aVktlhNhIsU_Lwpjn2_mfaqtD0ZX22odRhrtwk4LBztLqrySPvsolyrAfv4ntz2gMc5LGAe1ke3yfVZc9VTW6Gfe_WMtiA2J5tJoogOnm9nd4qVlR9SHn0zAxb_semi5rnufTkSfXS-E43qJwdvGyqj7KvpkUD3z5SRI1x2od9BwXrKCEFtCgNnOZY7lHSXn9xjh82mY0hNmpgO6npv1Qm7u4txZEg-nFaZ_SXj2GMP2KoOOQEpLKXzOjRlSA7oyTGXbQ7Q75a9cll_A4YgMRvtO-kR7OOFRhhjc7vUOwqg1ety5CPleshLt4umBG_mKhkpRmuomR06oNndyochWjEG_K7XEwLyR5bJmV-RTQUC2ka-2wgXVJDVk95cWCxYSV7PNyE_Jm9VWTL162maxwjDthTDWnADTJ-vkapRmqnd4UfLkHCWmsoPQN6XSQGmr4cghpL6N_KTW4ZPe6OIOtE9tJS6SLk092xPTu3CXl-U95mXs1ThdV7qf_-Pgq-JGF_vARoUhAhXAGT3i8AUO60Bta8I1gVohSfLxdaIDXbf3RKInr81cTz--GR72BmCBCUNm9DwO-UfCynVktt-CcaLCYfnAceRvWyeLi9-TT3Hle-r3yUAAXsb7h1an266ETpA7QMQODHRA9MSOE14KZeyziwlOC74T9aJreL_f5mVZTT2YjRJSiyTGk1eYXijDBGm6rrFhtYjHE8ZDi4LvLVh-Lx8CQSoz2kdkQ-WRa096KatZ5h02ufRcp3-Cn659FDNTQlhqyIjQueyyvkM1rvDgfbmo8hjEns-nRG0pDAGKiO2xND-GhL4RALIVQWWGH-lsJhDHfyrTOq6LMCodx-vQMTeywj2YfXW5dmenpPCrZf5KucO0kL1cP1Lh_4H-TPmZdVsIG8N_3h6G6WrKMbVYh4LgRRIGQs_H9Nx6EAc6cFsc_TFuLUsMGKgFBsAfLt_UJgQhd_fxvVLh_9QRLeBKz176qQSGEXaoad5odb7Aom7I1WIfGAA57-rdo2L8lQSRJAF0kdbiCsZeH6O2_GDTXom6yg2qs4Ah1F0G6xW-1Cs0WK6zpsNM7GSK0sk4277IXbYnS7f4nf9qwzgXmG3gBasssz7xgegZkZOtNrQAUGLR_ClYy7uE34igly0rsVJ80lgJy5QdKc8_DTlE3KrMnhVKNMw2L3VXfqTYWdumMpGIqcitnAP9uDk6WIFpjtS5G0mTLpPnf5J7BS3LItp2YfBX3saPudaWZpCVGXwOv9pd2ZVADAA34KvM_Rh9upF0n32RBO4MJk8VeHqLZc0_1qMQ3q6nVZoS6dpgFe2clzyWV3vkHfu9PfVUG2quD4MqslRzxJB9hplKxr-U3dIP6t_lycWGzOZhe1H6OlepntUxw5FEbI0lRl0bG7W64pIWQBxjh3o3jXrKT48Ppfg47CyUdhNLBWXpjnZVBwUnCmuBouoiFdddFX.Zd7-LSpoavpZ-jRL2K2BFg; _T_ANO=hmQME8/huXIVX5jj1oX2CK/hgEdlfOz7K6Bzhqk0iI0yXZByZkG8hOPxJ1I0/ucJy35sR/boZGNzTdXH4IY6kJ0O7S7RcPnB3xpgkiypAj/cuPPWxKLsDLGvd+7nP1ycCrUblMXPBsCosxdONO9/8S8T2aDuYU4Yf/xIzYQbkxpBeV61Oh30LTZg3G2flxUUdCm86L6aRx6keMVk3hdtdAbitlc7hWeaYPJ4DGRdNVJ487Z1BtNZDp1xXQjeoYbVWVSyFC75JEIPQ3Nq2wvDtHGZiKsdfexo6biHoRyMFNpK75tJTo0GzLlrFIo611VcscbtkRjP8bzmLYwqBPnHrA==" // rút gọn
)

// ============================================================================
// QUERY (dùng payload bạn cung cấp — giữ nguyên các type Long!)
// ============================================================================

const ViewerQuery = `
query viewerInfo($seriesId: Long!, $productId: Long!) {
  viewerInfo(seriesId: $seriesId, productId: $productId) {
    item {
      ...SingleFragment
    }
    seriesItem {
      ...SeriesFragment
    }
    prevItem {
      ...NearItemFragment
    }
    nextItem {
      ...NearItemFragment
    }
    viewerData {
      ...TextViewerData
      ...TalkViewerData
      ...ImageViewerData
      ...VodViewerData
    }
    displayAd {
      ...DisplayAd
    }
  }
}

fragment SingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  series {
    ...SeriesFragment
  }
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  discountRate
  discountRateText
  isShortsDrama
}

fragment SeriesFragment on Series {
  id
  seriesId
  title
  thumbnail
  landThumbnail
  categoryUid
  lang
  category
  categoryType
  subcategoryUid
  subcategory
  badge
  isAllFree
  isWaitfree
  ageGrade
  state
  onIssue
  authors
  description
  pubPeriod
  freeSlideCount
  lastSlideAddedDate
  waitfreeBlockCount
  waitfreePeriodByMinute
  bm
  saleState
  startSaleDt
  saleMethod
  discountRate
  discountRateText
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  translateProperty {
    ...TranslatePropertyFragment
  }
}

fragment ServicePropertyFragment on ServiceProperty {
  viewCount
  readCount
  ratingCount
  ratingSum
  commentCount
  pageContinue {
    ...ContinueInfoFragment
  }
  todayGift {
    ...TodayGift
  }
  preview {
    ...PreviewFragment
  }
  waitfreeTicket {
    ...WaitfreeTicketFragment
  }
  isAlarmOn
  isLikeOn
  ticketCount
  purchasedDate
  lastViewInfo {
    ...LastViewInfoFragment
  }
  purchaseInfo {
    ...PurchaseInfoFragment
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment ContinueInfoFragment on ContinueInfo {
  title
  isFree
  productId
  lastReadProductId
  scheme
  continueProductType
  hasNewSingle
  hasUnreadSingle
}

fragment TodayGift on TodayGift {
  id
  uid
  ticketType
  ticketKind
  ticketCount
  ticketExpireAt
  ticketExpiredText
  isReceived
  seriesId
}

fragment PreviewFragment on Preview {
  item {
    ...PreviewSingleFragment
  }
  nextItem {
    ...PreviewSingleFragment
  }
  usingScroll
}

fragment PreviewSingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
}

fragment OperatorPropertyFragment on OperatorProperty {
  thumbnail
  copy
  helixImpId
  isTextViewer
  selfCensorship
  isBook
  cashInfo {
    discountRate
    setDiscountRate
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment AssetPropertyFragment on AssetProperty {
  bannerImage
  cardImage
  cardTextImage
  cleanImage
  ipxVideo
  bannerSet {
    ...BannerSetFragment
  }
  cardSet {
    ...CardSetFragment
  }
  cardCover {
    ...CardCoverFragment
  }
}

fragment BannerSetFragment on BannerSet {
  backgroundImage
  backgroundColor
  mainImage
  titleImage
}

fragment CardSetFragment on CardSet {
  backgroundColor
  backgroundImage
}

fragment CardCoverFragment on CardCover {
  coverImg
  coverRestricted
}

fragment WaitfreeTicketFragment on WaitfreeTicket {
  chargedPeriod
  chargedCount
  chargedAt
}

fragment LastViewInfoFragment on LastViewInfo {
  isDone
  lastViewDate
  rate
  spineIndex
}

fragment PurchaseInfoFragment on PurchaseInfo {
  purchaseType
  rentExpireDate
  expired
}

fragment TranslatePropertyFragment on TranslateProperty {
  category {
    ...LocaleMapFragment
  }
  sub_category {
    ...LocaleMapFragment
  }
}

fragment LocaleMapFragment on LocaleMap {
  ko
  en
  th
}

fragment NearItemFragment on NearItem {
  productId
  slideType
  ageGrade
  isFree
  title
  thumbnail
}

fragment TextViewerData on TextViewerData {
  type
  atsServerUrl
  metaSecureUrl
  contentsList {
    chapterId
    contentId
    secureUrl
  }
}

fragment TalkViewerData on TalkViewerData {
  type
  talkDownloadData {
    dec
    host
    path
    talkViewerType
  }
}

fragment ImageViewerData on ImageViewerData {
  type
  imageDownloadData {
    ...ImageDownloadData
  }
}

fragment ImageDownloadData on ImageDownloadData {
  files {
    ...ImageDownloadFile
  }
  totalCount
  totalSize
  viewDirection
  gapBetweenImages
  readType
}

fragment ImageDownloadFile on ImageDownloadFile {
  no
  size
  secureUrl
  width
  height
}

fragment VodViewerData on VodViewerData {
  type
  vodDownloadData {
    contentPackId
    drmType
    endpointUrl
    width
    height
    duration
  }
  drmInfo {
    type
    serverType
    error
    fairplayCertificateUrl
    widevineLicenseUrl
    fairplayLicenseUrl
    token
    provider
    assertion
  }
}

fragment DisplayAd on DisplayAd {
  sectionUid
  bannerUid
  treviUid
  momentUid
}
`

// ============================================================================
// MODELS (mình chỉ khai báo những field cần dùng)
// ============================================================================

type GraphQLPayload struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

type GQLResponse struct {
	Data struct {
		ViewerInfo ViewerInfo `json:"viewerInfo"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type ViewerInfo struct {
	Item       Item       `json:"item"`
	SeriesItem *Series    `json:"seriesItem,omitempty"`
	PrevItem   *NearItem  `json:"prevItem,omitempty"`
	NextItem   *NearItem  `json:"nextItem,omitempty"`
	ViewerData ViewerData `json:"viewerData"`
}

type Item struct {
	Title string `json:"title"`
	// thêm field nếu cần về sau
	ProductId int `json:"productId,omitempty"`
	SeriesId  int `json:"seriesId,omitempty"`
}

type Series struct {
	SeriesId int    `json:"seriesId"`
	Title    string `json:"title"`
}

type NearItem struct {
	ProductId int    `json:"productId"`
	Title     string `json:"title"`
	// slideType, etc. omitted — thêm nếu cần
}

type ViewerData struct {
	AtsServerUrl string        `json:"atsServerUrl"`
	MetaSecure   string        `json:"metaSecureUrl,omitempty"`
	ContentsList []ContentPart `json:"contentsList"`
}

type ContentPart struct {
	ChapterId int    `json:"chapterId,omitempty"`
	ContentId int    `json:"contentId,omitempty"`
	SecureUrl string `json:"secureUrl"`
}

type ContentFileResponse struct {
	ContentInfo struct {
		ParagraphList []ParagraphNode `json:"paragraphList"`
	} `json:"contentInfo"`
}

type ParagraphNode struct {
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Children []ParagraphNode `json:"childParagraphList,omitempty"`
}

// ============================================================================
// MAIN: loop theo nextItem.productId — giữ flow tải/chạy y nguyên code gốc
// ============================================================================

func main() {
	// context tổng cho toàn bộ chạy
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}
	headers := map[string]string{
		"accept":          "*/*",
		"accept-language": "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7",
		"content-type":    "application/json",
		"origin":          "https://page.kakao.com",
		"referer":         "https://page.kakao.com/",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
		"cookie":          CookieString,
	}

	// Bắt đầu từ đây: set seriesId và productId chương ban đầu bạn muốn crawl
	currentSeries := 61097897
	currentProduct := 64177671

	for {
		fmt.Printf("\n>>> Crawl productId = %d (seriesId=%d)\n", currentProduct, currentSeries)

		// 1) build payload giữ nguyên như bạn yêu cầu
		payload := GraphQLPayload{
			OperationName: "viewerInfo",
			Query:         ViewerQuery,
			Variables: map[string]interface{}{
				"seriesId":  currentSeries,
				"productId": currentProduct,
			},
		}
		payloadBytes, _ := json.Marshal(payload)

		respBody, err := makeRequest(ctx, client, "POST", GraphQLEndpoint, headers, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Printf(">>> ERROR: GraphQL request failed: %v\n", err)
			break
		}

		var gqlResp GQLResponse
		if err := json.Unmarshal(respBody, &gqlResp); err != nil {
			fmt.Printf(">>> ERROR: Parse Error: %v\nBody: %s\n", err, string(respBody))
			break
		}
		if len(gqlResp.Errors) > 0 {
			fmt.Printf(">>> ERROR: GraphQL Error: %s\n", gqlResp.Errors[0].Message)
			break
		}

		viewerData := gqlResp.Data.ViewerInfo.ViewerData
		itemInfo := gqlResp.Data.ViewerInfo.Item
		nextItem := gqlResp.Data.ViewerInfo.NextItem

		if viewerData.AtsServerUrl == "" {
			fmt.Println(">>> ERROR: Không tìm thấy atsServerUrl. Kiểm tra lại Cookie.")
			break
		}

		// Tạo file
		safeTitle := sanitizeFilename(itemInfo.Title)
		fileName := safeTitle + ".txt"

		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf(">>> ERROR: Không thể tạo file %s: %v\n", fileName, err)
			break
		}
		fmt.Printf(">>> Đã tạo file: %s\n", fileName)
		fmt.Printf(">>> Tổng số phần: %d. Bắt đầu tải...\n", len(viewerData.ContentsList))

		// Ghi tiêu đề
		if _, err := file.WriteString(fmt.Sprintf("# %s\n\n", itemInfo.Title)); err != nil {
			fmt.Printf(">>> WARNING: lỗi ghi tiêu đề: %v\n", err)
		}

		// Loop tải các phần
		for i, part := range viewerData.ContentsList {
			if i > 0 {
				time.Sleep(1 * time.Second)
			}
			fullUrl := viewerData.AtsServerUrl + part.SecureUrl
			fmt.Printf(" -> [%d/%d] Đang tải... (Delay 1s)\n", i+1, len(viewerData.ContentsList))

			partBody, err := makeRequest(ctx, client, "GET", fullUrl, headers, nil)
			if err != nil {
				fmt.Printf("   [Lỗi tải phần %d]: %v\n", i+1, err)
				continue
			}

			var contentResp ContentFileResponse
			if err := json.Unmarshal(partBody, &contentResp); err != nil {
				// nếu không parse được JSON -> có thể là image / binary -> bỏ qua
				continue
			}

			for _, node := range contentResp.ContentInfo.ParagraphList {
				text := extractTextRecursive(node)
				text = strings.ReplaceAll(text, "&nbsp;", " ")
				if _, err := file.WriteString(text); err != nil {
					fmt.Printf("   Lỗi ghi file: %v\n", err)
				}
			}
		}

		file.Close()
		fmt.Println(">>> XONG! Kiểm tra file:", fileName)

		// Nếu không có nextItem -> dừng
		if nextItem == nil || nextItem.ProductId == 0 {
			fmt.Println(">>> HẾT CHƯƠNG hoặc không tìm thấy nextItem.")
			break
		}

		// Cập nhật productId để crawl chương tiếp theo
		fmt.Printf(">>> Next: %d (%s)\n", nextItem.ProductId, nextItem.Title)
		currentProduct = nextItem.ProductId

		// delay an toàn trước lần tiếp theo
		time.Sleep(2 * time.Second)
	}
}

// ============================================================================
// HELPERS
// ============================================================================

func makeRequest(ctx context.Context, client *http.Client, method, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status %d | Body: %s", resp.StatusCode, string(respBytes))
	}
	return respBytes, nil
}

func extractTextRecursive(node ParagraphNode) string {
	var sb strings.Builder
	if node.Text != "" {
		sb.WriteString(node.Text)
	}
	if node.Type == "BR" {
		sb.WriteString("\n")
	}
	for _, child := range node.Children {
		sb.WriteString(extractTextRecursive(child))
	}
	switch node.Type {
	case "P", "DIV":
		sb.WriteString("\n")
	case "H3":
		sb.WriteString("\n\n### ")
	}
	return sb.String()
}

func sanitizeFilename(name string) string {
	reg, _ := regexp.Compile(`[\\/:*?"<>|]`)
	safe := reg.ReplaceAllString(name, "_")
	return strings.TrimSpace(safe)
}
